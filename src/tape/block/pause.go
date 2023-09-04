package block

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Pause (silence) - ID 20
type Pause struct {
	pauseDuration int
}

func (p *Pause) Id() byte {
	return 0x20
}

func (p *Pause) Name() string {
	return "Pause (silence)"
}

func (p *Pause) Read(tzxFile *os.File) error {
	pauseDuration := make([]byte, 2)
	if _, err := tzxFile.Read(pauseDuration); err != nil {
		return err
	}
	p.pauseDuration = int(binary.LittleEndian.Uint16(pauseDuration))
	return nil
}

func (p *Pause) Info() string {
	return fmt.Sprintf("[duration: %dms]", p.pauseDuration)
}

func (p *Pause) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (p *Pause) PauseDuration() int {
	return p.pauseDuration
}
