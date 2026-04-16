//go:build !solution

package retryupdate

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

// getCurrentValue читает текущее значение ключа и возвращает новое значение
// после применения updateFn, а также версию, которую нужно передать в Set.
// Возвращает (newValue, oldVersion, err).
// При временной ошибке API возвращает ("", uuid.UUID{}, nil) — сигнал повторить.
func getCurrentValue(
	c kvapi.Client,
	key string,
	updateFn func(*string) (string, error),
) (newValue string, oldVersion uuid.UUID, err error) {
	rGet, errGet := c.Get(&kvapi.GetRequest{Key: key})

	var authErr *kvapi.AuthError
	var apiErr *kvapi.APIError

	switch {
	case errors.As(errGet, &authErr):
		return "", uuid.UUID{}, fmt.Errorf("error getting value: %w", errGet)

	case errors.Is(errGet, kvapi.ErrKeyNotFound):
		val, updateErr := updateFn(nil)
		if updateErr != nil {
			return "", uuid.UUID{}, fmt.Errorf("error in updateFn: %w", updateErr)
		}
		return val, uuid.UUID{}, nil

	case errors.As(errGet, &apiErr):
		// временная ошибка — пустой результат сигнализирует повторить GET
		return "", uuid.UUID{}, nil

	default:
		val, updateErr := updateFn(&rGet.Value)
		if updateErr != nil {
			return "", uuid.UUID{}, fmt.Errorf("error in updateFn: %w", updateErr)
		}
		return val, rGet.Version, nil
	}
}

// setResult описывает итог попытки записи.
type setResult int

const (
	setDone     setResult = iota // запись прошла успешно
	setRetry                     // временная ошибка, повторить Set с теми же параметрами
	setRefetch                   // конфликт версий, нужно заново читать значение
	setKeyGone                   // ключ пропал, нужно создать заново (updateFn(nil))
)

// trySet выполняет одну попытку Set и возвращает, что делать дальше.
func trySet(c kvapi.Client, req *kvapi.SetRequest) (setResult, error) {
	_, errSet := c.Set(req)

	var authErr *kvapi.AuthError
	var conflictErr *kvapi.ConflictError
	var apiErr *kvapi.APIError

	switch {
	case errSet == nil:
		return setDone, nil

	case errors.As(errSet, &authErr):
		return 0, fmt.Errorf("error setting value: %w", errSet)

	case errors.As(errSet, &conflictErr):
		// ExpectedVersion совпадает с тем, что мы отправили — наш предыдущий Set
		// на самом деле прошёл, просто вернул временную ошибку.
		if conflictErr.ExpectedVersion == req.NewVersion {
			return setDone, nil
		}
		return setRefetch, nil

	case errors.Is(errSet, kvapi.ErrKeyNotFound):
		return setKeyGone, nil

	case errors.As(errSet, &apiErr):
		return setRetry, nil

	default:
		return setDone, nil
	}
}

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	for {
		newValue, oldVersion, err := getCurrentValue(c, key, updateFn)
		if err != nil {
			return err
		}
		if newValue == "" && oldVersion == (uuid.UUID{}) {
			// временная ошибка GET — повторяем
			continue
		}

		req := &kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: oldVersion,
			NewVersion: uuid.Must(uuid.NewV4()),
		}

		for {
			result, err := trySet(c, req)
			if err != nil {
				return err
			}

			switch result {
			case setDone:
				return nil

			case setRetry:
				// повторяем Set с теми же параметрами

			case setRefetch:
				// чужой Set обогнал нас — читаем заново
				goto nextGet

			case setKeyGone:
				// ключ удалили — пересоздаём
				val, updateErr := updateFn(nil)
				if updateErr != nil {
					return fmt.Errorf("error in updateFn: %w", updateErr)
				}
				req.Value = val
				req.OldVersion = uuid.UUID{}
			}
		}

	nextGet:
	}
}
