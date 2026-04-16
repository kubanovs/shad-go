//go:build !solution

package ciletters

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func MakeLetter(n *Notification) (string, error) {
	tmpl, err := template.ParseFiles("templates/ciletter.tmpl")

	if err != nil {
		return "", nil
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, n)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	return buf.String(), nil
}

func (j Job) LastNLogLines(n int) []string {
	split := strings.Split(j.RunnerLog, "\n")
	if len(split) < n {
		return split
	}

	return split[len(split)-n:]
}
