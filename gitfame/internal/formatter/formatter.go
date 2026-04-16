package formatter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/internal/analyzer"
)

type personStatMarshall struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func Format(personStats []analyzer.PersonStat, orderBy string, format string) (string, error) {

	baseSortOrder := []string{"lines", "commits", "files", "name"}

	updatedOrder := reorderSortColumns(baseSortOrder, orderBy)

	comp := comparator{updatedOrder}

	sort.Slice(personStats, func(i, j int) bool {
		res := comp.sort(personStats[i], personStats[j])
		if res < 0 {
			return false
		}
		return true
	})

	switch format {
	case "tabular":
		return tabularFormat(personStats), nil
	case "csv":
		return csvFormat(personStats), nil
	case "json":
		return jsonFormat(personStats)
	case "json-lines":
		return jsonLinesFormat(personStats)
	default:
		return "", fmt.Errorf("unknown format to output: %s", format)
	}
}

func reorderSortColumns(sortOrder []string, orderBy string) []string {
	res := []string{orderBy}

	for _, el := range sortOrder {
		if el != orderBy {
			res = append(res, el)
		}
	}

	return res
}

func csvFormat(personStats []analyzer.PersonStat) string {
	var buf strings.Builder
	w := csv.NewWriter(&buf)

	// Записываем заголовок
	w.Write([]string{"Name", "Lines", "Commits", "Files"})

	// Записываем данные
	for _, stat := range personStats {
		w.Write([]string{
			stat.Name,
			strconv.Itoa(stat.Lines),
			strconv.Itoa(len(stat.Commits)),
			strconv.Itoa(len(stat.Files)),
		})
	}

	w.Flush()
	return buf.String()
}

func tabularFormat(personStats []analyzer.PersonStat) string {
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")

	for _, stat := range personStats {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\n", stat.Name, stat.Lines, len(stat.Commits), len(stat.Files))
	}

	w.Flush()

	return buf.String()
}

func jsonFormat(personStats []analyzer.PersonStat) (string, error) {
	var personStatsMarshall []personStatMarshall

	for _, stat := range personStats {
		personStatsMarshall = append(personStatsMarshall, personStatMarshall{
			stat.Name, stat.Lines,
			len(stat.Commits),
			len(stat.Files)})
	}

	json, err := json.Marshal(personStatsMarshall)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

func jsonLinesFormat(personStats []analyzer.PersonStat) (string, error) {
	var jsonLines []string

	for _, stat := range personStats {
		toMarshall := personStatMarshall{
			stat.Name, stat.Lines,
			len(stat.Commits),
			len(stat.Files)}

		jsonLine, err := json.Marshal(toMarshall)

		if err != nil {
			return "", err
		}
		jsonLines = append(jsonLines, string(jsonLine))
	}

	return strings.Join(jsonLines, "\n"), nil
}
