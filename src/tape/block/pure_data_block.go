package block

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
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

func (p *PureDataBlock) Info() [][]string {
	return [][]string{
		{"ZERO bit pulse length", strconv.Itoa(p.zeroBitPulseLength)},
		{"ONE bit pulse length", strconv.Itoa(p.oneBitPulseLength)},
		{"Used bits in last byte", strconv.Itoa(p.lastByteBitsUsed)},
		{"Pause after block", fmt.Sprintf("%d ms", p.pauseAfterBlock)},
		{"Data length", strconv.Itoa(p.dataSize)},
		{"Data flag byte", fmt.Sprintf("%x", p.dataFlag)},
	}
}

func (p *PureDataBlock) Pulses() []Pulse {
	pulses := make([]Pulse, 0)

	for _, dataByte := range p.data {
		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := p.zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = p.oneBitPulseLength
			}
			pulses = append(pulses, []Pulse{
				{
					Length: pulseLength,
					Level:  false,
				},
				{
					Length: pulseLength,
					Level:  true,
				},
			}...)
		}
	}

	return pulses
}

func (p *PureDataBlock) PauseDuration() int {
	return p.pauseAfterBlock
}
