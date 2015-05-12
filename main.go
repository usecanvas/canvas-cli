package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

var version = "0.1.1"
var usage = `
Usage:
	canvas new [<filename>] [--collection=COLLECTION]
	canvas list
	canvas pull <id> [--md | --json | --html]
	canvas delete <id>
	canvas account
	canvas login
	canvas env
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
	case args["new"].(bool):
		collection, path := decodeNew(args)
		cli.NewCanvas(collection, path)
	case args["list"].(bool):
		cli.ListCanvases("")
	case args["pull"].(bool):
		format := decodeFormat(args)
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

func decodeNew(args map[string]interface{}) (collection, filename string) {
	switch c := args["--collection"].(type) {
	case string:
		collection = c
	}

	switch f := args["<filename>"].(type) {
	case string:
		filename = f
	}
	return
}
