package main

import (
	"github.com/mark-rushakoff/ss33/cli"
)

func main() {
	app := cli.App()
	app.RunAndExitOnError()
}
