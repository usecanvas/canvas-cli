package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

var version = "0.2.0"
var usage = `
Usage:
	canvas account
	canvas delete <id>
	canvas env
	canvas list [--collection=COLLECTION]
	canvas login
	canvas new [<filename>] [--collection=COLLECTION]
	canvas pull <id> [--md | --json | --html]
	canvas search [<query>] [--collection=COLLECTION]
	canvas -h | --help
	canvas --version

Options:
  --md          Format Canvas as markdown
  --json        Format Canvas as json
  --html        Format Canvas as html
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
	case args["account"].(bool):
		cli.WhoAmI()
	case args["delete"].(bool):
		cli.DeleteCanvas(args["<id>"].(string))
	case args["env"].(bool):
		cli.Env()
	case args["list"].(bool):
		cli.ListCanvases(decodeCollection(args))
	case args["login"].(bool):
		cli.UserLogin()
		fmt.Println("Success!")
	case args["new"].(bool):
		cli.NewCanvas(decodeCollection(args), decodeFilename(args))
	case args["pull"].(bool):
		format := decodeFormat(args)
		cli.PullCanvas(args["<id>"].(string), format)
	case args["search"].(bool):
		if args["<query>"] != nil {
			query := args["<query>"].(string)
			cli.SearchUnix(decodeCollection(args), query)
		} else {
			cli.SearchInteractive(decodeCollection(args))
		}
	}
}

func decodeFormat(args map[string]interface{}) (format string) {
	switch {
	case args["--md"]:
		format = "md"
	case args["--json"]:
		format = "json"
	case args["--html"]:
		format = "html"
	default:
		format = "md"
	}
	return
}

func decodeCollection(args map[string]interface{}) (collection string) {
	switch c := args["--collection"].(type) {
	case string:
		collection = c
	}
	return
}

func decodeFilename(args map[string]interface{}) (filename string) {
	switch f := args["<filename>"].(type) {
	case string:
		filename = f
	}
	return
}

/*
#  display account info
 #  delete a canvas
 #  display CLI config
 #  list canvases in a collection
 #  login to canvas
 #  create a new canvas
 #  pull an existing canvas

 #  display this screen
 #  display version number
*/
