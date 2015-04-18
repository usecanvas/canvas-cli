package main

import "github.com/docopt/docopt-go"

var usage = `Canvas CLI

Usage:
	canvas account
	canvas new [<filename>]
	canvas -h | --help
	canvas --version

Options:
  -h, --help     Show this screen.
  --version      Show version.
`

func main() {
	arguments, _ := docopt.Parse(usage, nil, true, "Canvas CLI 0.1", false)
	cli := CLI{}

	switch {
	case arguments["new"].(bool):
		cli.New()
	case arguments["account"].(bool):
		cli.Account()
	}
}
