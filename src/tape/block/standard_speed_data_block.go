package block

import (
	"encoding/binary"
	"fmt"
	"os"
)

const StandardPilotToneLength = 2168
const StandardZeroBitPulseLength = 855
const StandardOneBitPulseLength = 1710

// StandardSpeedDataBlock - ID 10
type StandardSpeedDataBlock struct {
	pauseAfterBlock int
	dataSize        int
	dataFlag        byte
	data            []byte
}

func (s *StandardSpeedDataBlock) Id() byte {
	return 0x10
}

func (s *StandardSpeedDataBlock) Name() string {
	return "Standard Speed Data Block"
}

func (s *StandardSpeedDataBlock) Read(tzxFile *os.File) error {
	pauseAfterBlock := make([]byte, 2)
	if _, err := tzxFile.Read(pauseAfterBlock); err != nil {
		return err
	}
	s.pauseAfterBlock = int(binary.LittleEndian.Uint16(pauseAfterBlock))

	dataSize := make([]byte, 2)
	if _, err := tzxFile.Read(dataSize); err != nil {
		return err
	}
	s.dataSize = int(binary.LittleEndian.Uint16(dataSize))

	data := make([]byte, s.dataSize)
	if _, err := tzxFile.Read(data); err != nil {
		return err
	}
	s.data = data

	s.dataFlag = data[0]

	return nil
}

func (s *StandardSpeedDataBlock) Info() string {
	return fmt.Sprintf(
		"[data size: %d, flag: %x, after pause: %dms]",
		s.dataSize,
		s.dataFlag,
		s.pauseAfterBlock,
	)
}

func (s *StandardSpeedDataBlock) Pulses() []Pulse {
	pulses := make([]Pulse, 0)
	level := false

	// Generate pilot tone
	for i := 0; i < StandardPilotToneLength; i++ {
		pulses = append(pulses, Pulse{Length: StandardOneBitPulseLength, Level: level})
		level = !level
	}
	pulses = append(pulses, []Pulse{
		{
			Length: StandardZeroBitPulseLength,
			Level:  level,
		},
		{
			Length: StandardZeroBitPulseLength,
			Level:  !level,
		},
	}...)

	// Generate data
	for _, dataByte := range s.data {
		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := StandardZeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = StandardOneBitPulseLength
			}
			pulses = append(pulses, []Pulse{
				{
					Length: pulseLength,
					Level:  level,
				},
				{
					Length: pulseLength,
					Level:  !level,
				},
			}...)
		}
	}

	// Generate trailer
	for i := 0; i < 32; i++ {
		pulses = append(pulses, []Pulse{
			{
				Length: StandardOneBitPulseLength,
				Level:  level,
			},
			{
				Length: StandardOneBitPulseLength,
				Level:  !level,
			},
		}...)
	}

	return pulses
}

func (s *StandardSpeedDataBlock) PauseDuration() int {
	return s.pauseAfterBlock
}
