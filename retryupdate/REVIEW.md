# Code Review: retryupdate/update.go

## Итог

Задача решена корректно — все тесты проходят. Ниже описаны места, которые можно сделать чище и надёжнее.

---

## 1. Опечатка в сообщении об ошибке

**Строка 57** содержит текст `"error while get value"`, хотя речь идёт об ошибке метода `Set`:

```go
// было
return fmt.Errorf("error while get value: %w", errSet)

// должно быть
return fmt.Errorf("error setting value: %w", errSet)
```

---

## 2. Нарушение Go-конвенции: заглавная буква в тексте ошибки

**Строка 37**: в Go принято писать тексты ошибок строчными буквами без точки в конце
([https://go.dev/wiki/CodeReviewComments#error-strings](https://go.dev/wiki/CodeReviewComments#error-strings)).

```go
// было
return fmt.Errorf("Error while updateFn: %w", updateErr)

// должно быть
return fmt.Errorf("error in updateFn: %w", updateErr)
```

Аналогично остальные сообщения стоит привести к единому стилю: `"error in updateFn"`, `"error getting value"`, `"error setting value"`.

---

## 3. Избыточная инициализация нулевым значением

**Строка 20**: у `uuid.UUID` нулевое значение — это именно `uuid.UUID{}`, поэтому явное присвоение лишнее:

```go
// было
var oldVersion uuid.UUID = uuid.UUID{}

// достаточно
var oldVersion uuid.UUID
```

---

## 4. `default` без проверки `errGet == nil`

В `default`-ветке внешнего switch предполагается, что `errGet == nil` и `rGet` не nil. Если API неожиданно вернёт ошибку, не являющуюся `*APIError`, код обратится к `rGet.Value` при `rGet == nil` и упадёт с panic.

Надёжнее добавить явную проверку:

```go
case errGet == nil:
    val, updateErr := updateFn(&rGet.Value)
    ...
default:
    // неизвестная ошибка — можно залогировать или вернуть
    return errGet
```

---

## 5. Объявление переменных для ошибок лучше делать ближе к использованию

Сейчас `accessDenied` и `apiError` объявляются в начале внешнего цикла, хотя `accessDenied` переиспользуется и во внутреннем цикле. Это создаёт неявную связь между двумя циклами. Чище объявлять их непосредственно перед `errors.As`:

```go
var authErr *kvapi.AuthError
if errors.As(err, &authErr) {
    return fmt.Errorf("auth error: %w", err)
}
```

---

## 6. Дублирование вызовов `updateFn`

`updateFn` вызывается в трёх местах:

- при `ErrKeyNotFound` в GET
- при успешном GET (`default`)
- при `ErrKeyNotFound` в SET

Первые два случая можно объединить, выделив логику «получить текущее значение» в отдельный шаг перед циклом `Set`. Третий случай неизбежен, но его можно сделать явнее через `continue Outer` с переходом к GET вместо повторного вызова `updateFn` внутри цикла Set.

---

## 7. Сложность вложенных циклов с метками

Структура с двумя вложенными циклами и меткой `Outer` работает корректно и соответствует подсказке из README. Тем не менее, её можно упростить, если логику повторного чтения (GET + updateFn) вынести в отдельную вспомогательную функцию `getCurrentValue`, а основной цикл отвечал бы только за Set с ретраями.

---

## 8. Пример более чистой структуры

```go
func UpdateValue(c kvapi.Client, key string, updateFn func(*string) (string, error)) error {
    for {
        // --- шаг 1: читаем текущее значение ---
        var (
            oldVersion uuid.UUID
            newValue   string
        )

        rGet, err := c.Get(&kvapi.GetRequest{Key: key})
        switch {
        case isAuthError(err):
            return fmt.Errorf("auth error on get: %w", err)
        case errors.Is(err, kvapi.ErrKeyNotFound):
            v, err := updateFn(nil)
            if err != nil {
                return fmt.Errorf("error in updateFn: %w", err)
            }
            newValue = v
            // oldVersion остаётся нулевым
        case err != nil:
            continue // временная ошибка — ретраим
        default:
            v, err := updateFn(&rGet.Value)
            if err != nil {
                return fmt.Errorf("error in updateFn: %w", err)
            }
            newValue = v
            oldVersion = rGet.Version
        }

        // --- шаг 2: пишем новое значение ---
        newVersion := uuid.Must(uuid.NewV4())
        for {
            _, err := c.Set(&kvapi.SetRequest{
                Key:        key,
                Value:      newValue,
                OldVersion: oldVersion,
                NewVersion: newVersion,
            })
            switch {
            case err == nil:
                return nil
            case isAuthError(err):
                return fmt.Errorf("auth error on set: %w", err)
            case isConflict(err, newVersion):
                return nil // наш Set всё же прошёл
            case errors.Is(err, kvapi.ErrKeyNotFound):
                // ключ удалили — начинаем с чистого листа
                break // выйти из внутреннего цикла через continue внешнего
            default:
                continue // временная ошибка — ретраим Set
            }
            break
        }
        // продолжаем внешний цикл: идём на GET заново
    }
}

func isAuthError(err error) bool {
    var e *kvapi.AuthError
    return errors.As(err, &e)
}

func isConflict(err error, sentVersion uuid.UUID) bool {
    var e *kvapi.ConflictError
    if errors.As(err, &e) {
        return e.ExpectedVersion == sentVersion
    }
    return false
}
```

> Примечание: в примере выше внутренний `break` выходит из `switch`, а не из цикла.
> Для выхода из внутреннего `for` при `ErrKeyNotFound` в реальном коде всё равно потребуется метка или флаг — структура с `Outer` из оригинального решения здесь полностью оправдана.

---

## Резюме замечаний


| #   | Серьёзность    | Описание                                      |
| --- | -------------- | --------------------------------------------- |
| 1   | Баг / опечатка | Неверное сообщение об ошибке в Set-ветке      |
| 2   | Стиль          | Заглавная буква в тексте ошибки               |
| 3   | Стиль          | Лишняя инициализация нулевым значением        |
| 4   | Надёжность     | `default` без проверки `errGet == nil`        |
| 5   | Читаемость     | Объявление переменных далеко от использования |
| 6   | Читаемость     | Дублирование вызовов `updateFn`               |
| 7   | Архитектура    | Вложенные циклы можно разбить на функции      |


