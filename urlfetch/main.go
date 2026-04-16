//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	argsWithoutProgram := os.Args[1:]

	for _, url := range argsWithoutProgram {
		r, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v", err)
			os.Exit(1)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v", url, err)
			os.Exit(1)
		}

		fmt.Printf("%s", body)
	}
}
