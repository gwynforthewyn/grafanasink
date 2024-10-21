package main

import (
	"os"

	"github.com/playtechnique/grinksync"
)

func main() {
	os.Exit(grinksync.Main(os.Args[1:], os.Stdout))
}
