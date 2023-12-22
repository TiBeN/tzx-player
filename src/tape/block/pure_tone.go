package block

import (
	"encoding/binary"
	"os"
	"strconv"
)

// PureTone - ID 12
type PureTone struct {
	onePulseLength int
	pulsesNb       int
}

func (p *PureTone) Id() byte {
	return 0x12
}

func (p *PureTone) Name() string {
	return "Pure Tone"
}

func (p *PureTone) Read(tzxFile *os.File) error {
	onePulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(onePulseLength); err != nil {
		return err
	}
	p.onePulseLength = int(binary.LittleEndian.Uint16(onePulseLength))

	pulsesNb := make([]byte, 2)
	if _, err := tzxFile.Read(pulsesNb); err != nil {
		return err
	}
	p.pulsesNb = int(binary.LittleEndian.Uint16(pulsesNb))

	return nil
}

func (p *PureTone) Info() [][]string {
	return [][]string{
		{"One pulse length", strconv.Itoa(p.onePulseLength)},
		{"Number of pulses", strconv.Itoa(p.pulsesNb)},
	}
}

func (p *PureTone) Pulses() []Pulse {
	pulses := make([]Pulse, 0)
	level := false

	for i := 0; i < p.pulsesNb; i++ {
		pulses = append(pulses, Pulse{Length: p.onePulseLength, Level: level})
		level = !level
	}

	return pulses
}

func (p *PureTone) PauseDuration() int {
	return 0
}
