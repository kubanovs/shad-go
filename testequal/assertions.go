//go:build !solution

package testequal

import (
	"fmt"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if !areEqual(actual, expected) {
		t.Errorf("not equal:\nexpected: %v\nactual  : %v\nmessage : %s", expected, actual, msgRender(msgAndArgs...))
		return false
	}

	return true
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if areEqual(actual, expected) {
		t.Errorf("equal:\nexpected: %v\nactual  : %v\nmessage : %s", expected, actual, msgRender(msgAndArgs...))
		return false
	}

	return true
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if !areEqual(actual, expected) {
		t.Errorf("not equal:\nexpected: %v\nactual  : %v\nmessage : %s", expected, actual, msgRender(msgAndArgs...))
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if areEqual(actual, expected) {
		t.Errorf("equal:\nexpected: %v\nactual  : %v\nmessage : %s", expected, actual, msgRender(msgAndArgs...))
		t.FailNow()
	}
}

func msgRender(msgAndArgs ...interface{}) string {
	if msgAndArgs == nil || len(msgAndArgs) == 0 {
		return ""
	}

	if len(msgAndArgs) == 1 {
		msg, ok := msgAndArgs[0].(string)
		if !ok {
			return fmt.Sprintf("%+v", msgAndArgs[0])
		}
		return msg
	}

	return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
}

func areEqual(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}
	if expected == nil || actual == nil {
		return false
	}

	switch exp := expected.(type) {
	// Целые числа со знаком
	case int:
		act, ok := actual.(int)
		return ok && exp == act
	case int8:
		act, ok := actual.(int8)
		return ok && exp == act
	case int16:
		act, ok := actual.(int16)
		return ok && exp == act
	case int32:
		act, ok := actual.(int32)
		return ok && exp == act
	case int64:
		act, ok := actual.(int64)
		return ok && exp == act

	// Целые числа без знака
	case uint:
		act, ok := actual.(uint)
		return ok && exp == act
	case uint8:
		act, ok := actual.(uint8)
		return ok && exp == act
	case uint16:
		act, ok := actual.(uint16)
		return ok && exp == act
	case uint32:
		act, ok := actual.(uint32)
		return ok && exp == act
	case uint64:
		act, ok := actual.(uint64)
		return ok && exp == act

	// Строка
	case string:
		act, ok := actual.(string)
		return ok && exp == act

	// []int
	case []int:
		act, ok := actual.([]int)
		if !ok || len(exp) != len(act) {
			return false
		}

		if (exp == nil) != (act == nil) {
			return false
		}

		for i := range exp {
			if exp[i] != act[i] {
				return false
			}
		}
		return true

	// []byte
	case []byte:
		act, ok := actual.([]byte)
		if !ok || len(exp) != len(act) {
			return false
		}

		if (exp == nil) != (act == nil) {
			return false
		}

		for i := range exp {
			if exp[i] != act[i] {
				return false
			}
		}
		return true

	// map[string]string
	case map[string]string:
		act, ok := actual.(map[string]string)
		if !ok || len(exp) != len(act) {
			return false
		}

		if (exp == nil) != (act == nil) {
			return false
		}

		for key, val1 := range exp {
			val2, exists := act[key]
			if !exists || val1 != val2 {
				return false
			}
		}
		return true

	default:
		return false
	}
}
