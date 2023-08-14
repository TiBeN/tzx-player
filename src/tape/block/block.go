package block

import (
	"fmt"
	"os"
)

// Block holds information and content of a TZX tape data block
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
