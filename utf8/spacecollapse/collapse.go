//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
)

func CollapseSpaces(input string) string {
	if len(input) == 0 {
		return input
	}

	// Подсчитываем количество пробельных групп и непробельных символов
	var wsGroupCount int
	var countNonWs int
	prevWasSpace := false

	for _, r := range input {
		isSpace := unicode.IsSpace(r)

		if isSpace {
			if !prevWasSpace {
				wsGroupCount++
				prevWasSpace = true
			}
		} else {
			countNonWs++
			prevWasSpace = false
		}
	}

	var b strings.Builder
	// Аллоцируем память под результирующую строку
	b.Grow(countNonWs + wsGroupCount)

	prevWasSpace = false
	for _, r := range input {
		isSpace := unicode.IsSpace(r)

		if isSpace {
			if !prevWasSpace {
				b.WriteRune(' ')
				prevWasSpace = true
			}
		} else {
			b.WriteRune(r)
			prevWasSpace = false
		}
	}

	return b.String()
}
