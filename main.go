package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

var usage = `Canvas CLI

Usage:
	canvas account
	canvas login
	canvas new [<filename>]
	canvas pull <id>
	canvas list
	canvas -h | --help
	canvas --version

Options:
  -h, --help     Show this screen.
  --version      Show version.
`

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

//parses and validates argugmetns, then
//calls the appropriate method in the CLI class
//with the user submitted arguments
func main() {
	args, _ := docopt.Parse(usage, nil, true, "Canvas CLI 0.1", false)
	cli := NewCLI()

	switch {
	case args["new"].(bool):
		cli.NewCanvas()
	case args["account"].(bool):
		cli.WhoAmI()
	case args["list"].(bool):
		cli.ListCanvases()
	case args["login"].(bool):
		cli.Login()
	case args["pull"].(bool):
		cli.PullCanvas(args["<id>"].(string))
	}
}
