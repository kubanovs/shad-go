package formatter

import (
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/analyzer"
)

type comparator struct {
	columnOrder []string
}

func (c *comparator) sort(A analyzer.PersonStat, B analyzer.PersonStat) int {
	for _, el := range c.columnOrder {
		var res int

		switch el {
		case "lines":
			res = A.Lines - B.Lines
		case "commits":
			res = len(A.Commits) - len(B.Commits)
		case "files":
			res = len(A.Files) - len(B.Files)
		case "name":
			aNameLow := strings.ToLower(A.Name)
			bNameLow := strings.ToLower(B.Name)

			if aNameLow < bNameLow {
				res = 1
			} else {
				res = -1
			}
		}

		if res != 0 {
			return res
		}
	}

	return 1
}
