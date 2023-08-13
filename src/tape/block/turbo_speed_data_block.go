package block

import (
	"encoding/binary"
	"fmt"
	"os"
)

// TurboSpeedDataBlock - ID 11
type TurboSpeedDataBlock struct {
	pilotPulseLength      int
	pilotToneLength       int
	syncFirstPulseLength  int
	syncSecondPulseLength int
	zeroBitPulseLength    int
	oneBitPulseLength     int
	lastByteBitsUsed      int
	pauseAfterBlock       int
	dataSize              int
	dataFlag              byte
	data                  []byte
}

func (t *TurboSpeedDataBlock) Id() byte {
	return 0x11
}

func (t *TurboSpeedDataBlock) Name() string {
	return "Turbo Speed Data Block"
}

func (t *TurboSpeedDataBlock) Read(tzxFile *os.File) error {
	pilotPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(pilotPulseLength); err != nil {
		return err
	}
	t.pilotPulseLength = int(binary.LittleEndian.Uint16(pilotPulseLength))

	syncFirstPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(syncFirstPulseLength); err != nil {
		return err
	}
	t.syncFirstPulseLength = int(binary.LittleEndian.Uint16(syncFirstPulseLength))

	syncSecondPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(syncSecondPulseLength); err != nil {
		return err
	}
	t.syncSecondPulseLength = int(binary.LittleEndian.Uint16(syncSecondPulseLength))

	zeroBitPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(zeroBitPulseLength); err != nil {
		return err
	}
	t.zeroBitPulseLength = int(binary.LittleEndian.Uint16(zeroBitPulseLength))

	oneBitPulseLength := make([]byte, 2)
	if _, err := tzxFile.Read(oneBitPulseLength); err != nil {
		return err
	}
	t.oneBitPulseLength = int(binary.LittleEndian.Uint16(oneBitPulseLength))

	pilotToneLength := make([]byte, 2)
	if _, err := tzxFile.Read(pilotToneLength); err != nil {
		return err
	}
	t.pilotToneLength = int(binary.LittleEndian.Uint16(pilotToneLength))

	lastByteBitsUsed := make([]byte, 1)
	if _, err := tzxFile.Read(lastByteBitsUsed); err != nil {
		return err
	}
	t.lastByteBitsUsed = int(lastByteBitsUsed[0])

	pauseAfterBlock := make([]byte, 2)
	if _, err := tzxFile.Read(pauseAfterBlock); err != nil {
		return err
	}
	t.pauseAfterBlock = int(binary.LittleEndian.Uint16(pauseAfterBlock))

	dataSize := make([]byte, 3)
	if _, err := tzxFile.Read(dataSize); err != nil {
		return err
	}
	t.dataSize = int(binary.LittleEndian.Uint32(append(dataSize, 0)))

	data := make([]byte, t.dataSize)
	if _, err := tzxFile.Read(data); err != nil {
		return err
	}
	t.data = data

	t.dataFlag = data[0]

	return nil
}

func (t *TurboSpeedDataBlock) Info() string {
	return fmt.Sprintf(
		"[pilot p.: %dt, pilot tone length: %d, sync 1 p.: %dt, sync 2 p.: %dt, bit 0 p.: %dt, bit 1 p.: %dt, last byte used: %d, data size: %d, flag: %x, after pause: %dms]",
		t.pilotPulseLength,
		t.pilotToneLength,
		t.syncFirstPulseLength,
		t.syncSecondPulseLength,
		t.zeroBitPulseLength,
		t.oneBitPulseLength,
		t.lastByteBitsUsed,
		t.dataSize,
		t.dataFlag,
		t.pauseAfterBlock,
	)
}

func (t *TurboSpeedDataBlock) Samples(sampleRate int, bitDepth int) []byte {
	samples := make([]byte, 0)

	// Generate pilot
	lowLevel := true
	for i := 0; i < t.pilotToneLength; i++ {
		pulseSamples := GeneratePulseSamples(t.pilotPulseLength, sampleRate, bitDepth, lowLevel)
		samples = append(samples, pulseSamples...)
		lowLevel = !lowLevel
	}
	samples = append(samples, GeneratePulseSamples(t.zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel
	samples = append(samples, GeneratePulseSamples(t.zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel

	// Generate data
	for _, dataByte := range t.data {

		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := t.zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = t.oneBitPulseLength
			}
			samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
			samples = append(samples, GeneratePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
		}
	}

	// Generate trailer
	for i := 0; i < 32; i++ {
		samples = append(samples, GeneratePulseSamples(t.oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
		samples = append(samples, GeneratePulseSamples(t.oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
	}

	// Generate after block pause
	samples = append(samples, GeneratePause(t.pauseAfterBlock, sampleRate, bitDepth)...)

	return samples
}
