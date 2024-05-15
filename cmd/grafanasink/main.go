package main

import (
	"os"

	"github.com/playtechnique/grafanasink"
)

func main() {
	os.Exit(grafanasink.Main(os.Args[1:], os.Stdout))
}
