package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {

	var b strings.Builder
	b.Grow(len(input))

	for i := 0; i < len(input); {
		r, size := utf8.DecodeLastRuneInString(input[:len(input)-i])

		if r == utf8.RuneError {
			b.WriteRune(utf8.RuneError)
			i++
		} else {
			b.WriteRune(r)
			i += size
		}
	}

	return b.String()
}
