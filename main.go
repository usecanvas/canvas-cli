package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

var usage = `Canvas CLI 0.0.1
Usage:
	canvas new [<filename>]
	canvas list [--collection]
	canvas pull <id> [-f | --format=<format>]
	canvas account
	canvas login
	canvas -h | --help
	canvas --version

Options:
  -h, --help             Show this screen.
  --version              Show version.
  --collection           Document collection (defaults to current user).
  -f, --format <format>  Format: md, json, or git
`

//unified error handler
func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

//parses and validates argugments, then
//calls the appropriate method in the CLI class
//with the user submitted arguments
func main() {
	args, _ := docopt.Parse(usage, nil, true, "Canvas CLI 0.1", false)
	cli := NewCLI()

	switch {
	case args["new"].(bool):
		switch path := args["<filename>"].(type) {
		case string:
			cli.NewCanvas(path)
		case nil:
			cli.BlankCanvas()
		}
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
