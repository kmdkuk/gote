package main

import (
	"github.com/kmdkuk/gote/cmd"
)

// These variables are set in build step
var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	cmd.Execute()
}
