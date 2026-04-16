//go:build !solution

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	overrider *overrider
	stack     Stack
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{NewOverrider(), Stack{}}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	split := strings.Split(row, " ")

	if isOverriding(split) {
		return e.stack.slice, e.override(split)
	} else {
		expanded, err := e.toBaseOperations(split)

		if err != nil {
			return e.stack.slice, err
		}

		for _, op := range expanded {
			switch {
			case isNumber(op):
				num, _ := strconv.Atoi(op)
				e.stack.push(num)

			case op == "+":
				if e.stack.size < 2 {
					return e.stack.slice, fmt.Errorf("not enough arguments for +")
				}
				a := e.stack.pop()
				b := e.stack.pop()
				e.stack.push(b + a)

			case op == "-":
				if e.stack.size < 2 {
					return e.stack.slice, fmt.Errorf("not enough arguments for -")
				}
				a := e.stack.pop()
				b := e.stack.pop()
				e.stack.push(b - a)

			case op == "*":
				if e.stack.size < 2 {
					return e.stack.slice, fmt.Errorf("not enough arguments for *")
				}
				a := e.stack.pop()
				b := e.stack.pop()
				e.stack.push(b * a)

			case op == "/":
				if e.stack.size < 2 {
					return e.stack.slice, fmt.Errorf("not enough arguments for /")
				}
				a := e.stack.pop()
				b := e.stack.pop()
				if a == 0 {
					return e.stack.slice, fmt.Errorf("division by zero")
				}
				e.stack.push(b / a)

			case op == "dup":
				if e.stack.size < 1 {
					return e.stack.slice, fmt.Errorf("not enough arguments for dup")
				}
				val := e.stack.peek()
				e.stack.push(val)

			case op == "drop":
				if e.stack.size < 1 {
					return e.stack.slice, fmt.Errorf("not enough arguments for drop")
				}
				e.stack.pop()

			case op == "swap":
				if e.stack.size < 2 {
					return e.stack.slice, fmt.Errorf("not enough arguments for swap")
				}
				a := e.stack.pop()
				b := e.stack.pop()
				e.stack.push(a)
				e.stack.push(b)

			case op == "over":
				if e.stack.size < 2 {
					return e.stack.slice, fmt.Errorf("not enough arguments for over")
				}
				second := e.stack.slice[e.stack.size-2]
				e.stack.push(second)
			}
		}

		return e.stack.slice, nil
	}
}

func isOverriding(args []string) bool {
	return len(args) >= 4 && args[0] == ":" && args[len(args)-1] == ";"
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (e *Evaluator) override(args []string) error {
	if _, err := strconv.Atoi(args[1]); err == nil {
		return fmt.Errorf("Can't set new operation: %v", err)
	}

	definition := args[1]
	description := strings.Join(args[2:len(args)-1], " ")

	if err := e.overrider.SetOverride(definition, description); err != nil {
		return fmt.Errorf("Can't set new operation: %v", err)
	}

	return nil
}

func (e *Evaluator) toBaseOperations(args []string) ([]string, error) {
	var updated []string
	for _, el := range args {
		lowerEl := strings.ToLower(el)

		if isNumber(lowerEl) {
			updated = append(updated, lowerEl)
			continue
		}

		val, err := e.overrider.GetOverride(el, false)
		if err != nil {
			return nil, errors.New("not existed definition")
		}
		updated = append(updated, strings.Split(val, " ")...)
	}
	return updated, nil
}
