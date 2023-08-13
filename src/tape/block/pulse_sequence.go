package block

import (
	"encoding/binary"
	"fmt"
	"os"
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

func (p *PulseSequence) Info() string {
	return fmt.Sprintf("[pulses  nb: %d]", p.pulsesNb)
}

func (p *PulseSequence) Samples(sampleRate int, bitDepth int) []byte {
	samples := make([]byte, 0)
	lowLevel := true

	for i := 0; i < len(p.pulsesLengths); i = i + 2 {
		pulseLength := int(binary.LittleEndian.Uint16(p.pulsesLengths[i : i+2]))
		samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
	}

	return samples
}
