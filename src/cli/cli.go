package cli

import (
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"os"
)

const ConvertDefaultSamplingRate = 44100
const ConvertDefaultBitDepth = 8
const ConvertDefaultSpeedFactor = 1.0

type Cli struct {
	tapeService *tape.Service
	commands    []Command
}

func NewCli(tapeService *tape.Service) *Cli {
	return &Cli{
		tapeService: tapeService,
		commands: []Command{
			&Convert{},
			&Info{},
			&Play{},
		},
	}
}

// Exec parses command argument then execute matching command
func (c *Cli) Exec(args []string) {
	if len(args) < 2 {
		c.err(fmt.Errorf("no command specified"))
	}

	if args[1] == "help" {
		c.usage()
		os.Exit(0)
	}

	var cmd Command
	for _, j := range c.commands {
		if j.Name() == args[1] {
			cmd = j
			break
		}
	}

	if cmd == nil {
		c.err(fmt.Errorf("unknown command '%s'", args[1]))
	}

	if err := cmd.Exec(c.tapeService, args[2:]); err != nil {
		c.err(fmt.Errorf("%s: %s", cmd.Name(), err.Error()))
	}
}

func (c *Cli) err(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	c.usage()
	os.Exit(3)
}

func (c *Cli) usage() {
	fmt.Println("")
	fmt.Println("Usage: tzx-player COMMAND [cmd opts]")
	fmt.Println("")
	fmt.Println("Play 8bits computers data tapes as TZX files")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Printf("  %-20sShow this help message\n", "help")
	for _, cmd := range c.commands {
		fmt.Printf("  %-20s%s\n", cmd.Name(), cmd.Description())
		fmt.Print(cmd.Usage())
	}
}
