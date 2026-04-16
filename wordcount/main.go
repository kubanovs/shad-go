//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	argsWithoutProg := os.Args[1:]

	mp := make(map[string]int)

	var scanner bufio.Scanner

	for i := 0; i < len(argsWithoutProg); i++ {
		f, err := os.Open(argsWithoutProg[i])
		check(err)
		scanner = *bufio.NewScanner(f)
		for scanner.Scan() {
			l := scanner.Text()
			mp[l]++
		}
	}

	for k, v := range mp {
		if v > 1 {
			fmt.Printf("%d\t%s\n", v, k)
		}
	}

}
