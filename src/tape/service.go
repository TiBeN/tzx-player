package tape

import (
	"io"
	"time"
)

// Service is the main interface for actions related to a TZX Tape
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ConvertToWavFile converts the given TZX tape file into an audio PCM WAV file
func (s *Service) ConvertToWavFile(tzxFile string, outputFile string, samplingRate int, bitDepth int) (*time.Duration, error) {
	start := time.Now()

	tape, err := NewTape(tzxFile)
	if err != nil {
		return nil, err
	}

	tapeReader, err := NewReader(tape, samplingRate, bitDepth)
	if err != nil {
		return nil, err
	}

	wavWriter, err := NewWavFileWriter(outputFile, samplingRate, bitDepth)
	defer func() {
		err = wavWriter.Close()
	}()

	if err != nil {
		return nil, err
	}
	buf := make([]byte, 65535)
	if _, err = io.CopyBuffer(wavWriter, tapeReader, buf); err != nil {
		return nil, err
	}

	end := time.Since(start)
	return &end, nil
}

func (s *Service) Info(tzxFile string) ([]string, error) {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return nil, err
	}
	return tape.Info(), nil
}
