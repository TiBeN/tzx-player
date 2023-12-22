package block

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

// DirectRecording - ID 15
type DirectRecording struct {
	nbTstatePerSample int
	pauseAfterBlock   int
	lastByteBitsUsed  int
	samplesDataSize   int
	samplesData       []byte
}

func (d *DirectRecording) Id() byte {
	return 0x11
}

func (d *DirectRecording) Name() string {
	return "Direct Recording"
}

func (d *DirectRecording) Read(tzxFile *os.File) error {
	var nbTstatePerSample uint16
	if err := binary.Read(tzxFile, binary.LittleEndian, &nbTstatePerSample); err != nil {
		return err
	}
	d.nbTstatePerSample = int(nbTstatePerSample)

	var pauseAfterBlock uint16
	if err := binary.Read(tzxFile, binary.LittleEndian, &pauseAfterBlock); err != nil {
		return err
	}
	d.pauseAfterBlock = int(pauseAfterBlock)

	var lastByteBitsUsed uint8
	if err := binary.Read(tzxFile, binary.LittleEndian, &lastByteBitsUsed); err != nil {
		return err
	}
	d.lastByteBitsUsed = int(lastByteBitsUsed)

	samplesDataSize := make([]byte, 3)
	if _, err := tzxFile.Read(samplesDataSize); err != nil {
		return err
	}
	d.samplesDataSize = int(binary.LittleEndian.Uint32(append(samplesDataSize, 0)))

	samplesData := make([]byte, d.samplesDataSize)
	if _, err := tzxFile.Read(samplesData); err != nil {
		return err
	}
	d.samplesData = samplesData

	return nil
}

func (d *DirectRecording) Info() [][]string {
	return [][]string{
		{"Number of T-states per sample", strconv.Itoa(d.nbTstatePerSample)},
		{"Pause after block", fmt.Sprintf("%d ms", d.pauseAfterBlock)},
		{"Used bits in last byte", strconv.Itoa(d.lastByteBitsUsed)},
		{"Samples 'data length", strconv.Itoa(d.samplesDataSize)},
	}
}

func (d *DirectRecording) Pulses() []Pulse {
	pulses := make([]Pulse, 0)

	currentPulse := Pulse{}
	for _, samples := range d.samplesData {
		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			level := int(samples)&i > 0
			if currentPulse.Level != level {
				if currentPulse.Length > 0 {
					pulses = append(pulses, currentPulse)
				}
				currentPulse = Pulse{Level: level}
			}
			currentPulse.Length += d.nbTstatePerSample
		}
	}
	// Store last pulse
	if currentPulse.Length > 0 {
		pulses = append(pulses, currentPulse)
	}

	return pulses
}

func (d *DirectRecording) PauseDuration() int {
	return d.pauseAfterBlock
}
