package main

import (
	"errors"
	"strings"
)

type overrider struct {
	mp             map[string]string
	baseOperations map[string]struct{}
}

func NewOverrider() *overrider {
	return &overrider{make(map[string]string), map[string]struct{}{
		"+":    {},
		"-":    {},
		"*":    {},
		"/":    {},
		"dup":  {},
		"over": {},
		"drop": {},
		"swap": {},
	}}
}

func (obj *overrider) GetOverride(definition string, baseOpSensetive bool) (string, error) {
	definition = strings.ToLower(definition)
	val, ok := obj.mp[definition]
	_, isBaseOperation := obj.baseOperations[definition]

	if isBaseOperation && (baseOpSensetive || !ok) {
		return definition, nil
	}

	if !ok {
		return "", errors.New("Not existed definition")
	}

	return val, nil
}

func (obj *overrider) SetOverride(definition string, description string) error {
	definition = strings.ToLower(definition)
	split := strings.Split(description, " ")
	var updated []string

	for _, el := range split {
		lowerEl := strings.ToLower(el)

		if isNumber(lowerEl) {
			updated = append(updated, lowerEl)
			continue
		}

		val, err := obj.GetOverride(strings.ToLower(el), false)
		if err != nil {
			return errors.New("Not existed definition")
		}
		updated = append(updated, val)
	}

	obj.mp[definition] = strings.Join(updated, " ")
	return nil
}
