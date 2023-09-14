package block

import (
	"encoding/binary"
	"fmt"
	"os"

	"golang.org/x/text/encoding/charmap"
)

var TextIds map[byte]string

func init() {
	TextIds = map[byte]string{
		0x00: "Full title",
		0x01: "Software house/publisher",
		0x02: "Author(s)",
		0x03: "Year of publication",
		0x04: "Language",
		0x05: "Game/utility type",
		0x06: "Price",
		0x07: "Protection scheme/loader",
		0x08: "Origin",
		0xFF: "Comment(s)",
	}
}

type ArchiveInfo struct {
	texts []Text
}

type Text struct {
	textId byte
	text   string
}

func (a *ArchiveInfo) Id() byte {
	return 0x32
}

func (a *ArchiveInfo) Name() string {
	return "Archive Info"
}

func (a *ArchiveInfo) Read(tzxFile *os.File) error {
	blockLength := make([]byte, 2)
	if _, err := tzxFile.Read(blockLength); err != nil {
		return err
	}
	var textNb uint8
	if err := binary.Read(tzxFile, binary.LittleEndian, &textNb); err != nil {
		return err
	}

	fmt.Println(textNb)

	for i := 0; i < int(textNb); i++ {
		text := Text{}

		textId := make([]byte, 1)
		if _, err := tzxFile.Read(textId); err != nil {
			return err
		}
		text.textId = textId[0]

		var textLen uint8
		if err := binary.Read(tzxFile, binary.LittleEndian, &textLen); err != nil {
			return err
		}

		textBytes := make([]byte, textLen)
		if _, err := tzxFile.Read(textBytes); err != nil {
			return err
		}

		// Strings are encoded using ISO charset
		isoDecoder := charmap.ISO8859_1.NewDecoder()
		decodedBytes, err := isoDecoder.Bytes(textBytes)
		if err != nil {
			return err
		}

		text.text = string(decodedBytes)
		a.texts = append(a.texts, text)
	}

	return nil
}

func (a *ArchiveInfo) Info() string {
	info := "["

	for i, text := range a.texts {
		info += fmt.Sprintf("id: %s, text: %s", TextIds[text.textId], text.text)
		if i < len(a.texts)-1 {
			info += ", "
		}
	}

	return info + "]"
}

func (a *ArchiveInfo) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (a *ArchiveInfo) PauseDuration() int {
	return 0
}
