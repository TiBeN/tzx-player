package tape

import (
	"encoding/binary"
	"os"
)

const HeaderLength = 44

// WavFileWriter is an io.Writer which writes raw audio PCM data to a WAV file.
// Calling WaveFileWriter.Close() after writing data is mandatory to generate a valid
// Wav file (Header data is written at this time).
type WavFileWriter struct {
	f          *os.File
	SampleRate int
	BitDepth   int
	dataLength int
}

func NewWavFileWriter(fileName string, sampleRate int, BitDepth int) (*WavFileWriter, error) {
	w := &WavFileWriter{
		SampleRate: sampleRate,
		BitDepth:   BitDepth,
	}

	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	w.f = f

	// Seek to data part, after header
	if _, err = w.f.Seek(HeaderLength, 0); err != nil {
		return nil, err
	}

	return w, nil
}

// Write raw samples to file
func (w *WavFileWriter) Write(samples []byte) (n int, err error) {
	w.dataLength += len(samples)
	return w.f.Write(samples)
}

// Close writes Wav header data then close the file
func (w *WavFileWriter) Close() (err error) {
	defer func() {
		err = w.f.Close()
	}()

	// Seek to beginning of the file
	if _, err := w.f.Seek(0, 0); err != nil {
		return err
	}

	// Write header
	header := w.GenerateHeader()
	if _, err := w.f.Write(header); err != nil {
		return err
	}

	return nil
}

func (w *WavFileWriter) GenerateHeader() []byte {
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")

	fileSize := (uint32)(w.dataLength + HeaderLength - 8)
	binary.LittleEndian.PutUint32(header[4:8], fileSize)

	copy(header[8:16], "WAVEfmt ")

	blocSize := []byte{0x10, 0x00, 0x00, 0x00}
	copy(header[16:20], blocSize)

	audioFormat := []byte{0x01, 0x00}
	copy(header[20:22], audioFormat)

	channelNb := []byte{0x01, 0x00}
	copy(header[22:24], channelNb)

	binary.LittleEndian.PutUint32(header[24:28], uint32(w.SampleRate))

	bytePerSec := uint32(w.SampleRate * (w.BitDepth / 8))
	binary.LittleEndian.PutUint32(header[28:32], bytePerSec)

	bytePerBloc := uint16(w.BitDepth / 8)
	binary.LittleEndian.PutUint16(header[32:34], bytePerBloc)

	binary.LittleEndian.PutUint16(header[34:36], uint16(w.BitDepth))

	copy(header[36:40], "data")

	dataSize := uint32(w.dataLength)
	binary.LittleEndian.PutUint32(header[40:44], dataSize)

	return header
}
