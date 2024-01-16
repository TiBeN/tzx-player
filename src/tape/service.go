package tape

import (
	"io"
	"time"
)

// Service is the main entry point for actions related to a TZX Tape
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ConvertToWavFile converts the given TZX tape file into an audio PCM WAV file
func (s *Service) ConvertToWavFile(tzxFile string, outputFile string, samplingRate int, bitDepth int, speedFactor float64) (*time.Duration, error) {
	start := time.Now()

	tape, err := NewTape(tzxFile)
	if err != nil {
		return nil, err
	}

	tapeReader, err := NewReader(tape, samplingRate, bitDepth, speedFactor)
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

// Info returns information about a TZX tape file (version, blocks etc.)
func (s *Service) Info(tzxFile string) (*TapeInfo, error) {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return nil, err
	}
	info := tape.Info()
	return &info, nil
}

// Play plays a TZX file through audio sound card
func (s *Service) Play(tzxFile string, samplingRate int, bitDepth int, speedFactor float64) (*Player, error) {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return nil, err
	}

	tapeReader, err := NewReader(tape, samplingRate, bitDepth, speedFactor)
	if err != nil {
		return nil, err
	}

	player := NewPlayer(tapeReader)
	if err = player.Start(); err != nil {
		return nil, err
	}

	return player, nil
}
