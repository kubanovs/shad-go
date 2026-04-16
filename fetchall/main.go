//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type UrlResponse struct {
	Size int
	Secs float64
	Url  string
	Err  error
}

func fetchUrl(url string, ch chan UrlResponse) {
	start := time.Now()
	r, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v", err)
		ch <- UrlResponse{Err: err}
		return
	}
	f, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %s: %v", url, err)
		ch <- UrlResponse{Err: err}
	}
	ch <- UrlResponse{len(f), time.Since(start).Seconds(), url, nil}
}

func main() {
	argsWithoutProgram := os.Args[1:]
	ch := make(chan UrlResponse)
	start := time.Now()
	for _, url := range argsWithoutProgram {
		go fetchUrl(url, ch)
	}

	for range argsWithoutProgram {
		resp := <-ch
		if resp.Err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", resp.Err)
		} else {
			fmt.Printf("%.2fs %7d %s\n", resp.Secs, resp.Size, resp.Url)
		}
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
