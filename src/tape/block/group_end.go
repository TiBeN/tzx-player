package block

import "os"

// GroupEnd - ID 22
type GroupEnd struct {
}

func (g *GroupEnd) Id() byte {
	return 0x22
}

func (g *GroupEnd) Name() string {
	return "Group end"
}

func (g *GroupEnd) Read(tzxFile *os.File) error {
	return nil
}

func (g *GroupEnd) Info() string {
	return ""
}

func (g *GroupEnd) Samples(sampleRate int, bitDepth int) []byte {
	return make([]byte, 0)
}
