package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const TzxSignature = "ZXTape!"

type Tape struct {
	Header Header
	Blocks []Block
}

type Header struct {
	MajorVersion int
	MinorVersion int
}

func NewTape(tzxFile string) (*Tape, error) {
	f, err := os.Open(tzxFile)
	if err != nil {
		return nil, err
	}

	tape := Tape{}
	if err := tape.readHeader(f); err != nil {
		return nil, err
	}

	if err := tape.readBlocks(f); err != nil {
		return nil, err
	}

	// read blocks
	return &tape, nil
}

// Info returns information about the tape
func (t *Tape) Info() []string {
	infos := make([]string, len(t.Blocks)+1)
	infos[0] = fmt.Sprintf("TZX tape, version %d.%d", t.Header.MajorVersion, t.Header.MinorVersion)
	for blockNb, block := range t.Blocks {
		infos[blockNb+1] = fmt.Sprintf("Block %d: %s %s", blockNb, block.Name(), block.Info())
	}
	return infos
}

func (t *Tape) Samples(sampleRate int, bitDepth int) []byte {
	samples := make([]byte, 0)
	for _, block := range t.Blocks {
		samples = append(samples, block.Samples(sampleRate, bitDepth)...)
	}
	return samples
}

// Read header data from TZX file
func (t *Tape) readHeader(tzxFile *os.File) error {
	headerBytes := make([]byte, 10)
	if _, err := tzxFile.Read(headerBytes); err != nil {
		return err
	}

	if string(headerBytes[0:7]) != TzxSignature {
		return errors.New("not a valid TZX file (no TZX signature in header)")
	}
	if headerBytes[7] != 0x1a {
		return errors.New("not a valid TZX file (end of text file marker not found in header)")
	}
	t.Header = Header{
		MajorVersion: int(headerBytes[8]),
		MinorVersion: int(headerBytes[9]),
	}
	return nil
}

// Read blocks from TZX file
func (t *Tape) readBlocks(tzxFile *os.File) error {
	for {
		blockId := make([]byte, 1)
		if _, err := tzxFile.Read(blockId); err == io.EOF {
			break
		}
		block, err := NewBlock(blockId[0], tzxFile)
		if err != nil {
			return err
		}
		t.Blocks = append(t.Blocks, block)
	}
	return nil
}
