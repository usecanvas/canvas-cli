package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

var version = "0.0.2"
var usage = `
Usage:
	canvas new [<filename>]
	canvas list
	canvas pull <id>
	canvas delete <id>
	canvas account
	canvas login
	canvas env
	canvas -h | --help
	canvas --version

Options:
  -h, --help    Show this screen.
  --version     Show version.
`

//TODO:
// --collection  Document collection (defaults to current user).
// [--format=<format>]
//  --format      Format: md, json, or git

//unified error handler
func check(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

//parses and validates argugments, then
//calls the appropriate method in the CLI class
//with the user submitted arguments
func main() {
	args, _ := docopt.Parse(usage, nil, true, "Canvas CLI "+version, false)
	if len(args) == 0 {
		check(errors.New("Could not parse command line options"))
	}

	cli := NewCLI()
	switch {
	case args["new"].(bool):
		switch path := args["<filename>"].(type) {
		case string:
			cli.NewCanvasPath(path)
		case nil:
			cli.NewCanvas()
		}
	case args["list"].(bool):
		cli.ListCanvases("")
	case args["pull"].(bool):
		cli.PullCanvas(args["<id>"].(string))
	case args["delete"].(bool):
		cli.DeleteCanvas(args["<id>"].(string))
	case args["account"].(bool):
		cli.WhoAmI()
	case args["login"].(bool):
		cli.Login()
		fmt.Println("Success!")
	case args["env"].(bool):
		cli.Env()
	}
}
