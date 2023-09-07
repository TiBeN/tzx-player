package block

import (
	"fmt"
	"os"
)

// MessageBlock - ID 31
type MessageBlock struct {
	displayDuration int
	message         string
}

func (m *MessageBlock) Id() byte {
	return 0x31
}

func (m *MessageBlock) Name() string {
	return "Message Block"
}

func (m *MessageBlock) Read(tzxFile *os.File) error {
	displayDuration := make([]byte, 1)
	if _, err := tzxFile.Read(displayDuration); err != nil {
		return err
	}
	m.displayDuration = int(displayDuration[0])

	messageLength := make([]byte, 1)
	if _, err := tzxFile.Read(messageLength); err != nil {
		return err
	}

	message := make([]byte, messageLength[0])
	if _, err := tzxFile.Read(message); err != nil {
		return err
	}
	m.message = string(message)

	return nil
}

func (m *MessageBlock) Info() string {
	return fmt.Sprintf("[display duration: %d, message: %s]", m.displayDuration, m.message)
}

func (m *MessageBlock) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (m *MessageBlock) PauseDuration() int {
	return 0
}
