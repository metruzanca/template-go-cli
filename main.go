package main

import "github.com/metruzanca/go-cli/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute()
}
