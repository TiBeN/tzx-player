package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Block interface {
	// Id returns the identifier of this block
	Id() byte

	// Name returns the name of this block
	Name() string

	// Read block data from the given TZX file.
	// It is expected the offset of the file descriptor is positioned
	// at the beginning of the block data (just after the block ID byte)
	Read(tzxFile *os.File) error

	// Info returns information about this block
	Info() string

	// Samples Generates and returns audio PCM samples
	// for this block
	Samples(sampleRate int, bitDepth int) []byte
}

func NewBlock(id byte, tzxFile *os.File) (Block, error) {
	var block Block

	switch id {
	case 0x10:
		block = &StandardSpeedDataBlock{}
	case 0x11:
		block = &TurboSpeedDataBlock{}
	case 0x12:
		block = &PureTone{}
	case 0x13:
		block = &PulseSequence{}
	case 0x14:
		block = &PureDataBlock{}
	case 0x20:
		block = &Pause{}
	case 0x21:
		block = &GroupStart{}
	case 0x22:
		block = &GroupEnd{}
	case 0x30:
		block = &TextDescription{}
	default:
		return nil, fmt.Errorf("unknown block id %x", id)
	}

	if err := block.Read(tzxFile); err != nil {
		return nil, err
	}
	return block, nil
}

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
		pulseSamples := generatePulseSamples(oneBitPulseLength, sampleRate, bitDepth, lowLevel)
		samples = append(samples, pulseSamples...)
		lowLevel = !lowLevel
	}
	samples = append(samples, generatePulseSamples(zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel
	samples = append(samples, generatePulseSamples(zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel

	// Generate data
	for _, dataByte := range s.data {

		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = oneBitPulseLength
			}
			samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
			samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
		}
	}

	// Generate trailer
	for i := 0; i < 32; i++ {
		samples = append(samples, generatePulseSamples(oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
		samples = append(samples, generatePulseSamples(oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
	}

	// Generate after block pause
	samples = append(samples, generatePause(s.pauseAfterBlock, sampleRate, bitDepth)...)

	return samples
}

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
		pulseSamples := generatePulseSamples(t.pilotPulseLength, sampleRate, bitDepth, lowLevel)
		samples = append(samples, pulseSamples...)
		lowLevel = !lowLevel
	}
	samples = append(samples, generatePulseSamples(t.zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel
	samples = append(samples, generatePulseSamples(t.zeroBitPulseLength, sampleRate, bitDepth, lowLevel)...)
	lowLevel = !lowLevel

	// Generate data
	for _, dataByte := range t.data {

		for i := 128; i >= 1; i = i / 2 { // Iterate over every bit
			pulseLength := t.zeroBitPulseLength
			if int(dataByte)&i > 0 {
				pulseLength = t.oneBitPulseLength
			}
			samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
			samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
		}
	}

	// Generate trailer
	for i := 0; i < 32; i++ {
		samples = append(samples, generatePulseSamples(t.oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
		samples = append(samples, generatePulseSamples(t.oneBitPulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
	}

	// Generate after block pause
	samples = append(samples, generatePause(t.pauseAfterBlock, sampleRate, bitDepth)...)

	return samples
}

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
		pulseSamples := generatePulseSamples(p.onePulseLength, sampleRate, bitDepth, lowLevel)
		samples = append(samples, pulseSamples...)
		lowLevel = !lowLevel
	}

	return samples
}

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
		samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
		lowLevel = !lowLevel
	}

	return samples
}

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
			samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
			samples = append(samples, generatePulseSamples(pulseLength, sampleRate, bitDepth, lowLevel)...)
			lowLevel = !lowLevel
		}
	}

	// Generate after block pause
	samples = append(samples, generatePause(p.pauseAfterBlock, sampleRate, bitDepth)...)

	return samples
}

// Pause (silence) - ID 20
type Pause struct {
	pauseDuration int
}

func (t *Pause) Id() byte {
	return 0x20
}

func (p *Pause) Name() string {
	return "Pause (silence)"
}

func (p *Pause) Read(tzxFile *os.File) error {
	pauseDuration := make([]byte, 2)
	if _, err := tzxFile.Read(pauseDuration); err != nil {
		return err
	}
	p.pauseDuration = int(binary.LittleEndian.Uint16(pauseDuration))
	return nil
}

func (p *Pause) Info() string {
	return fmt.Sprintf("[duration: %dms]", p.pauseDuration)
}

func (p *Pause) Samples(sampleRate int, bitDepth int) []byte {
	return generatePause(p.pauseDuration, sampleRate, bitDepth)
}

// GroupStart - ID 21
type GroupStart struct {
	nameLength int
	name       string
}

func (g *GroupStart) Id() byte {
	return 0x21
}

func (g *GroupStart) Name() string {
	return "Group start"
}

func (g *GroupStart) Read(tzxFile *os.File) error {
	nameLength := make([]byte, 1)
	if _, err := tzxFile.Read(nameLength); err != nil {
		return err
	}
	g.nameLength = int(nameLength[0])

	name := make([]byte, g.nameLength)
	if _, err := tzxFile.Read(name); err != nil {
		return err
	}
	g.name = string(name)

	return nil
}

func (g *GroupStart) Info() string {
	return fmt.Sprintf("[name: %s, length: %d]", g.name, g.nameLength)
}

func (g *GroupStart) Samples(sampleRate int, bitDepth int) []byte {
	return make([]byte, 0)
}

// GroupEnd - ID 22
type GroupEnd struct {
}

func (g *GroupEnd) Id() byte {
	return 0x22
}

func (g *GroupEnd) Name() string {
	return "Group end"
}

func (g *GroupEnd) Read(tzxFile *os.File) error {
	return nil
}

func (g *GroupEnd) Info() string {
	return ""
}

func (g *GroupEnd) Samples(sampleRate int, bitDepth int) []byte {
	return make([]byte, 0)
}

// TextDescription - ID 30
type TextDescription struct {
	description string
}

func (t *TextDescription) Id() byte {
	return 0x30
}

func (t *TextDescription) Name() string {
	return "Text Description"
}

func (t *TextDescription) Read(tzxFile *os.File) error {
	textLength := make([]byte, 1)
	if _, err := tzxFile.Read(textLength); err != nil {
		return err
	}
	textDescription := make([]byte, textLength[0])
	if _, err := tzxFile.Read(textDescription); err != nil {
		return err
	}
	t.description = string(textDescription)
	return nil
}

func (t *TextDescription) Info() string {
	return fmt.Sprintf("[description: %s]", t.description)
}

func (t *TextDescription) Samples(sampleRate int, bitDepth int) []byte {
	return make([]byte, 0)
}
