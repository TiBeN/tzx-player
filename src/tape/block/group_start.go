package block

import (
	"os"
	"strconv"
)

// GroupStart - ID 21
type GroupStart struct {
	nameLength int
	name       string
}

func (g *GroupStart) Id() byte {
	return 0x21
}

func (g *GroupStart) Name() string {
	return "Group start"
}

func (g *GroupStart) Read(tzxFile *os.File) error {
	nameLength := make([]byte, 1)
	if _, err := tzxFile.Read(nameLength); err != nil {
		return err
	}
	g.nameLength = int(nameLength[0])

	name := make([]byte, g.nameLength)
	if _, err := tzxFile.Read(name); err != nil {
		return err
	}
	g.name = string(name)

	return nil
}

func (g *GroupStart) Info() [][]string {
	return [][]string{
		{"Group name string length", strconv.Itoa(g.nameLength)},
		{"Group name", g.name},
	}
}

func (g *GroupStart) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (g *GroupStart) PauseDuration() int {
	return 0
}
