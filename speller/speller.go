//go:build !solution

package speller

var (
	units = map[int64]string{
		0:  "zero",
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
		11: "eleven",
		12: "twelve",
		13: "thirteen",
		14: "fourteen",
		15: "fifteen",
		16: "sixteen",
		17: "seventeen",
		18: "eighteen",
		19: "nineteen",
	}

	tens = map[int64]string{
		20: "twenty",
		30: "thirty",
		40: "forty",
		50: "fifty",
		60: "sixty",
		70: "seventy",
		80: "eighty",
		90: "ninety",
	}

	scales = []struct {
		value int64
		name  string
	}{
		{1000000000000, "trillion"},
		{1000000000, "billion"},
		{1000000, "million"},
		{1000, "thousand"},
		{100, "hundred"},
	}
)

// Spell converts a number to its American English word representation
func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}

	if n < 0 {
		return "minus " + Spell(-n)
	}

	return spellPositive(n)
}

// spellPositive handles positive numbers only
func spellPositive(n int64) string {
	if n == 0 {
		return ""
	}

	// Handle numbers less than 100
	if n < 100 {
		return spellUnder100(n)
	}

	// Handle larger numbers using scales
	for _, scale := range scales {
		if n >= scale.value {
			quotient := n / scale.value
			remainder := n % scale.value

			if remainder == 0 {
				return spellPositive(quotient) + " " + scale.name
			}
			return spellPositive(quotient) + " " + scale.name + " " + spellPositive(remainder)
		}
	}

	return ""
}

// spellUnder100 converts numbers from 1 to 99
func spellUnder100(n int64) string {
	if n <= 19 {
		return units[n]
	}

	if n%10 == 0 {
		return tens[n]
	}

	return tens[(n/10)*10] + "-" + units[n%10]
}
