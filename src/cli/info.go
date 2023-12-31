package cli

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
)

type Info struct {
}

func (c *Info) Name() string {
	return "info"
}

func (c *Info) Description() string {
	return "Output TZX tape informations"
}

func (c *Info) Usage() string {
	usage := fmt.Sprintf("    Args:\n")
	usage += fmt.Sprintf("      tzx-player info INPUT_TZX_FILE\n")
	return usage
}

func (c *Info) Exec(service *tape.Service, args []string) error {
	if len(args) < 1 {
		return errors.New("TZX file not specified")
	}

	info, err := service.Info(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("%-40s: %s\n", "TZX Tape Version", info.Version)

	for _, block := range info.Blocks {
		fmt.Println("")
		for _, params := range block {
			fmt.Printf("%-40s: %s\n", params[0], params[1])
		}
	}

	return nil
}
