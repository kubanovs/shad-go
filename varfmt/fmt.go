//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Sprintf(format string, args ...interface{}) string {
	var b strings.Builder
	var argsString []string
	var argMaxLen int
	var argsInFormatCount int

	// находим аргумент с максимальной длиной
	for _, arg := range args {
		argsString = append(argsString, fmt.Sprint(arg))
		argMaxLen = max(len(argsString[len(argsString)-1]), argMaxLen)
	}

	var isOpen bool
	var indexTextSize int

	// Считаем количество аргументов в форматируемой строке и количество индексов
	for _, r := range format {

		if r == '{' {
			isOpen = true
			argsInFormatCount++
		} else if r == '}' {
			isOpen = false
		} else if isOpen {
			indexTextSize++
		}
	}

	// Аллоцируем память
	b.Grow(len(format) - indexTextSize - 2*argsInFormatCount + argsInFormatCount*argMaxLen)
	isOpen = false
	var countArgsInFormat int
	var openIndex int
	// Генерируем результирующую строку
	for i := 0; i < len(format); {
		r, size := utf8.DecodeRuneInString(format[i:])

		if r == '{' {
			isOpen = true
			openIndex = i + 1
		} else if r == '}' {
			isOpen = false
			var ind int
			if openIndex != i {
				ind, _ = strconv.Atoi(format[openIndex:i])
			} else {
				ind = countArgsInFormat
			}

			countArgsInFormat++
			b.WriteString(argsString[ind])
		} else if !isOpen {
			b.WriteRune(r)
		}
		i += size
	}

	return b.String()
}
