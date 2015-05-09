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
	canvas new [<filename>] [--collection=COLLECTION]
	canvas list
	canvas pull <id> [--format=(md|html|json) [default: md]]
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
// canvas delete <id>

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
		// note: do don't know the Username until we've authed
		// and for now auth is the sole domain of the Cli struct
		var collection string
		switch c := args["--collection"].(type) {
		case string:
			collection = c
		}

		switch path := args["<filename>"].(type) {
		case string:
			cli.NewCanvasPath(collection, path)
		case nil:
			cli.NewCanvas(collection)
		}
	case args["list"].(bool):
		cli.ListCanvases("")
	case args["pull"].(bool):
		var format string
		switch f := args["--format"].(type) {
		case string:
			format = f
		case nil:
			format = "md"
		}

		cli.PullCanvas(args["<id>"].(string), format)
	case args["delete"].(bool):
		cli.DeleteCanvas(args["<id>"].(string))
	case args["account"].(bool):
		cli.WhoAmI()
	case args["login"].(bool):
		cli.UserLogin()
		fmt.Println("Success!")
	case args["env"].(bool):
		cli.Env()
	}
}
