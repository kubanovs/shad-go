//go:build !solution

package main

import (
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/gfcmd"
)

func main() {
	os.Exit(gfcmd.Main())
}
