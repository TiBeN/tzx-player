package tape

import (
	"bytes"
	"io"
)

// Service is the main interface for actions related to a TZX Tape
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ConvertToWavFile converts the given TZX tape file into an audio PCM WAV file
func (s *Service) ConvertToWavFile(tzxFile string, outputFile string, samplingRate int, bitDepth int) error {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return err
	}

	wavWriter, err := NewWavFileWriter(outputFile, samplingRate, bitDepth)
	if err != nil {
		return err
	}
	if _, err = io.Copy(wavWriter, bytes.NewReader(tape.Samples(samplingRate, bitDepth))); err != nil {
		return err
	}

	return wavWriter.Close()
}

func (s *Service) Info(tzxFile string) ([]string, error) {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return nil, err
	}
	return tape.Info(), nil
}
