package main

import (
	"github.com/TiBeN/tzx-player/cli"
	"github.com/TiBeN/tzx-player/tape"
	"os"
)

func main() {
	c := cli.NewCli(tape.NewService())
	c.Exec(os.Args)
}
