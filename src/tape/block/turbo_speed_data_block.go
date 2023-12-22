package block

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
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

func (t *TurboSpeedDataBlock) Info() [][]string {
	return [][]string{
		{"PILOT pulse length", strconv.Itoa(t.pilotPulseLength)},
		{"SYNC first pulse length", strconv.Itoa(t.syncFirstPulseLength)},
		{"SYNC second pulse length", strconv.Itoa(t.syncSecondPulseLength)},
		{"ZERO bit pulse length", strconv.Itoa(t.zeroBitPulseLength)},
		{"ONE bit pulse length", strconv.Itoa(t.oneBitPulseLength)},
		{"PILOT tone length", strconv.Itoa(t.pilotToneLength)},
		{"Used bits in last byte", strconv.Itoa(t.lastByteBitsUsed)},
		{"Pause after block", fmt.Sprintf("%d ms", t.pauseAfterBlock)},
		{"Data length", strconv.Itoa(t.dataSize)},
		{"Data flag byte", fmt.Sprintf("%x", t.dataFlag)},
	}
}

func (t *TurboSpeedDataBlock) Pulses() []Pulse {
	pulses := make([]Pulse, 0)
	level := false

	// Generate pilot tone
	for i := 0; i < t.pilotToneLength; i++ {
		pulses = append(pulses, Pulse{Length: t.pilotPulseLength, Level: level})
		level = !level
	}

	// Generate sync pulses
	pulses = append(pulses, []Pulse{
		{
			Length: t.syncFirstPulseLength,
			Level:  level,
		},
		{
			Length: t.syncSecondPulseLength,
			Level:  !level,
		},
	}...)

	// Generate data pulses
	for _, dataByte := range t.data {
		// @TODO: don't convert unused bits in last byte
		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := t.zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = t.oneBitPulseLength
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
				Length: t.oneBitPulseLength,
				Level:  level,
			},
			{
				Length: t.oneBitPulseLength,
				Level:  !level,
			},
		}...)
	}

	return pulses
}

func (t *TurboSpeedDataBlock) PauseDuration() int {
	return t.pauseAfterBlock
}
