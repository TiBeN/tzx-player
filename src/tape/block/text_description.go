package block

import (
	"os"
)

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

func (t *TextDescription) Info() [][]string {
	return [][]string{
		{"Description", t.description},
	}
}

func (t *TextDescription) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (t *TextDescription) PauseDuration() int {
	return 0
}
