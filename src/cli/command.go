package cli

import "github.com/TiBeN/tzx-player/tape"

// Command is the main interface for command line commands
type Command interface {

	// Name returns the name of the command
	Name() string

	// Description returns a short description of the command
	Description() string

	// Usage returns the documentation of the command
	Usage() string

	// Exec parses the cli args and executes the command
	Exec(service *tape.Service, args []string) error
}
