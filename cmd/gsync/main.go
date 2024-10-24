package main

import (
	"os"

	"github.com/playtechnique/gsync"
)

func main() {
	os.Exit(gsync.Main(os.Args[1:], os.Stdout))
}
