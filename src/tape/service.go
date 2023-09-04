package tape

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
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

func (s *Service) Play(tzxFile string, samplingRate int, bitDepth int) error {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return err
	}

	tapeReader, err := NewReader(tape, samplingRate, bitDepth)
	if err != nil {
		return err
	}

	samples, _ := io.ReadAll(tapeReader)
	samplesReader := bytes.NewReader(samples)

	if err := portaudio.Initialize(); err != nil {
		return err
	}

	buf := make([]byte, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(samplingRate), len(buf), &buf)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err = stream.Start(); err != nil {
		return err
	}
	defer stream.Stop()

	fmt.Printf("start playing %s ...\n", tzxFile)

	for remain := len(samples); remain > 0; remain -= len(buf) {
		if len(buf) > remain {
			buf = buf[:remain]
		}
		err = binary.Read(samplesReader, binary.BigEndian, buf)
		if err == io.EOF {
			break
		}

		stream.Write()
	}

	if err = portaudio.Terminate(); err != nil {
		return err
	}

	fmt.Println("stop")

	return nil
}
