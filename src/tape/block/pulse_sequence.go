package block

import (
	"encoding/binary"
	"os"
	"strconv"
)

// PulseSequence - ID 13
type PulseSequence struct {
	pulsesNb      int
	pulsesLengths []byte
}

func (p *PulseSequence) Id() byte {
	return 0x13
}

func (p *PulseSequence) Name() string {
	return "Pulse Sequence"
}

func (p *PulseSequence) Read(tzxFile *os.File) error {
	pulsesNb := make([]byte, 1)
	if _, err := tzxFile.Read(pulsesNb); err != nil {
		return err
	}
	p.pulsesNb = int(pulsesNb[0])

	p.pulsesLengths = make([]byte, p.pulsesNb*2)
	if _, err := tzxFile.Read(p.pulsesLengths); err != nil {
		return err
	}

	return nil
}

func (p *PulseSequence) Info() [][]string {
	return [][]string{
		{"Pulses number", strconv.Itoa(p.pulsesNb)},
	}
}

func (p *PulseSequence) Pulses() []Pulse {
	pulses := make([]Pulse, 0)
	level := false

	for i := 0; i < len(p.pulsesLengths); i = i + 2 {
		pulses = append(pulses, Pulse{
			Length: int(binary.LittleEndian.Uint16(p.pulsesLengths[i : i+2])),
			Level:  level,
		})
		level = !level
	}

	return pulses
}

func (p *PulseSequence) PauseDuration() int {
	return 0
}
