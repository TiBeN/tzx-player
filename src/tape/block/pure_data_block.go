package block

import (
	"encoding/binary"
	"fmt"
	"os"
)

// PureDataBlock - ID 14
type PureDataBlock struct {
	zeroBitPulseLength int
	oneBitPulseLength  int
	lastByteBitsUsed   int
	pauseAfterBlock    int
	dataSize           int
	dataFlag           byte
	data               []byte
}

func (p *PureDataBlock) Id() byte {
	return 0x14
}

func (p *PureDataBlock) Name() string {
	return "Pure Data Block"
}

func (p *PureDataBlock) Read(tzxFile *os.File) error {
	zeroBitPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(zeroBitPulseLength); err != nil {
		return err
	}
	p.zeroBitPulseLength = int(binary.LittleEndian.Uint16(zeroBitPulseLength))

	oneBitPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(oneBitPulseLength); err != nil {
		return err
	}
	p.oneBitPulseLength = int(binary.LittleEndian.Uint16(oneBitPulseLength))

	lastByteBitsUsed := make([]byte, 1)
	if _, err := tzxFile.Read(lastByteBitsUsed); err != nil {
		return err
	}
	p.lastByteBitsUsed = int(lastByteBitsUsed[0])

	pauseAfterBlock := make([]byte, 2)
	if _, err := tzxFile.Read(pauseAfterBlock); err != nil {
		return err
	}
	p.pauseAfterBlock = int(binary.LittleEndian.Uint16(pauseAfterBlock))

	dataSize := make([]byte, 3)
	if _, err := tzxFile.Read(dataSize); err != nil {
		return err
	}
	p.dataSize = int(binary.LittleEndian.Uint32(append(dataSize, 0)))

	data := make([]byte, p.dataSize)
	if _, err := tzxFile.Read(data); err != nil {
		return err
	}
	p.data = data

	return nil
}

func (p *PureDataBlock) Info() string {
	return fmt.Sprintf(
		"[bit 0 p.: %dt, bit 1 p.: %dt, last byte used: %d, data size: %d, flag: %x, after pause: %dms]",
		p.zeroBitPulseLength,
		p.oneBitPulseLength,
		p.lastByteBitsUsed,
		p.dataSize,
		p.dataFlag,
		p.pauseAfterBlock,
	)
}

func (p *PureDataBlock) Samples(sampleRate int, bitDepth int) []byte {
	samples := make([]byte, 0)
	lowLevel := true

	// Generate data
	for _, dataByte := range p.data {

		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := p.zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = p.oneBitPulseLength
			}
			samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
			samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
		}
	}

	// Generate after block pause
	samples = append(samples, GeneratePause(p.pauseAfterBlock, sampleRate, bitDepth)...)

	return samples
}
