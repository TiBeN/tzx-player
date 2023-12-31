package tape

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape/block"
	"io"
	"os"
	"strconv"
)

const TzxSignature = "ZXTape!"

type Tape struct {
	Header   Header
	Blocks   []block.Block
	FileName string
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

	tape := Tape{
		FileName: tzxFile,
	}
	if err := tape.readHeader(f); err != nil {
		return nil, err
	}

	if err := tape.readBlocks(f); err != nil {
		return nil, err
	}

	// read blocks
	return &tape, nil
}

type TapeInfo struct {
	Version string
	Blocks  [][][]string
}

// Info returns information about the tape
func (t *Tape) Info() TapeInfo {
	info := TapeInfo{
		Version: fmt.Sprintf("%d.%d", t.Header.MajorVersion, t.Header.MinorVersion),
	}

	for i, blk := range t.Blocks {
		blockInfo := [][]string{
			{"Block Number", strconv.Itoa(i + 1)},
			{"Block ID", fmt.Sprintf("%x", blk.Id())},
			{"Block Type", blk.Name()},
		}
		for _, param := range blk.Info() {
			blockInfo = append(blockInfo, param)
		}
		info.Blocks = append(info.Blocks, blockInfo)
	}

	return info
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
		b, err := block.NewBlock(blockId[0], tzxFile)
		if err != nil {
			return err
		}
		t.Blocks = append(t.Blocks, b)
	}
	return nil
}
