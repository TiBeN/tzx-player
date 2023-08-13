package block

import (
	"encoding/binary"
	"fmt"
	"os"
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

func (p *PureTone) Info() string {
	return fmt.Sprintf("[pulse length: %d, pulses number: %d]", p.onePulseLength, p.pulsesNb)
}

func (p *PureTone) Samples(sampleRate int, bitDepth int) []byte {
	samples := make([]byte, 0)

	// Generate pilot
	lowLevel := true
	for i := 0; i < p.pulsesNb; i++ {
		pulseSamples := GeneratePulseSamples(p.onePulseLength, sampleRate, bitDepth, lowLevel)
		samples = append(samples, pulseSamples...)
		lowLevel = !lowLevel
	}

	return samples
}
