# Canvas CLI

```bash
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
```

## Installation

OS X - Download the binary [using github releases](https://github.com/usecanvas/canvas-cli/releases).

## Usage

### `canvas new`

Create a new canvas and output the URL to STDOUT.

Examples:

Creating a blank document

    $ canvas new
    https://beta.usecanvas.com/csquared/-/d6ffa6aa-63d1-44b0-8d99-07d25a9db115

Creating from a file

    $ canvas new README.md
    https://beta.usecanvas.com/csquared/-/dbad8e34-0b2b-4ce7-80b6-7dc395b6a28e

Creating from STDIN

    $ cat README.md | canvas new
    https://beta.usecanvas.com/csquared/-/fb8f489a-ecba-47fa-840c-325cf13d7885

Create and open in a browser

    $ canvas new | xargs open

### `canvas pull`

Output document in format to STDOUT

    $ canvas pull fb8f489a-ecba-47fa-840c-325cf13d7885
    # Canvas CLI

    ...

Creating HTML

    $ npm install -g marked
    $ canvas pull fb8f489a-ecba-47fa-840c-325cf13d7885 | marked
    <h1 id="canvas-cli">Canvas CLI</h1>

    ...

**note**: this behavior will soon be replaced with the `--format` flag.


## Development

### Build

    ./build.sh

### Development

    go run *.go
