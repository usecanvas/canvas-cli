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
	arguments, _ := docopt.Parse(usage, nil, true, "Canvas CLI 0.1", false)
	cli := NewCLI()

	switch {
	case arguments["new"].(bool):
		cli.NewCanvas()
	case arguments["account"].(bool):
		cli.WhoAmI()
	case arguments["list"].(bool):
		cli.List()
	case arguments["login"].(bool):
		cli.Login()
	}
}
