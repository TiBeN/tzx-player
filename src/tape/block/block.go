package block

import (
	"fmt"
	"os"
)

// Block holds information and content of a TZX tape data block
// @TODO: Implements others blocks types
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

	// Pulses generates and returns the pulses of this block
	Pulses() []Pulse

	// PauseDuration return the trailing pause duration of this block in ms
	PauseDuration() int
}

type Pulse struct {
	// Length of the pulse in T state per second
	Length int

	// Low of High level (false = low, true = high)
	Level bool
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
	case 0x15:
		block = &DirectRecording{}
	case 0x20:
		block = &Pause{}
	case 0x21:
		block = &GroupStart{}
	case 0x22:
		block = &GroupEnd{}
	case 0x30:
		block = &TextDescription{}
	case 0x31:
		block = &MessageBlock{}
	case 0x32:
		block = &ArchiveInfo{}
	case 0x33:
		block = &HardwareType{}
	default:
		return nil, fmt.Errorf("unknown block id %x", id)
	}

	if err := block.Read(tzxFile); err != nil {
		return nil, err
	}
	return block, nil
}
