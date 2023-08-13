package block

import (
	"encoding/binary"
	"fmt"
	"os"
)

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

func (s *StandardSpeedDataBlock) Samples(sampleRate int, bitDepth int) []byte {
	samples := make([]byte, 0)
	pilotToneLength := 2168
	zeroBitPulseLength := 855
	oneBitPulseLength := 1710

	// Generate pilot
	lowLevel := true
	for i := 0; i < pilotToneLength; i++ {
		pulseSamples := GeneratePulseSamples(oneBitPulseLength, sampleRate, bitDepth, lowLevel)
		samples = append(samples, pulseSamples...)
		lowLevel = !lowLevel
	}
	samples = append(samples, GeneratePulseSamples(zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel
	samples = append(samples, GeneratePulseSamples(zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel

	// Generate data
	for _, dataByte := range s.data {

		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = oneBitPulseLength
			}
			samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
			samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
		}
	}

	// Generate trailer
	for i := 0; i < 32; i++ {
		samples = append(samples, GeneratePulseSamples(oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
		samples = append(samples, GeneratePulseSamples(oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
	}

	// Generate after block pause
	samples = append(samples, GeneratePause(s.pauseAfterBlock, sampleRate, bitDepth)...)

	return samples
}
