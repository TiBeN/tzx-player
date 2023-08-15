package tape

import (
	"encoding/binary"
	"os"
)

func writeToWavFile(filename string, samplerate int, bitdepth int, samples []byte) {
	f, _ := os.Create(filename)

	f.Write([]byte("RIFF"))

	length := (uint32)((len(samples) * 2) + 44 - 8) // data
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, length)
	f.Write(lengthBytes)

	f.Write([]byte("WAVEfmt "))
	f.Write([]byte{0x10, 0x00, 0x00, 0x00}) // BlocSize

	f.Write([]byte{0x01, 0x00}) // AudioFormat
	f.Write([]byte{0x01, 0x00}) // ChannelNb

	samplerateBytes := make([]byte, 4) // SampleRate
	binary.LittleEndian.PutUint32(samplerateBytes, uint32(samplerate))
	f.Write(samplerateBytes)

	bytePerSec := make([]byte, 4) // BytePerSec
	binary.LittleEndian.PutUint32(bytePerSec, uint32(samplerate*2))
	f.Write(bytePerSec)

	f.Write([]byte{0x02, 0x00}) // BytePerBloc
	f.Write([]byte{0x10, 0x00}) // BitsPerSample

	f.Write([]byte("data"))

	dataSize := uint32(len(samples) * 2) //DataSize
	dataSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSizeBytes, dataSize)
	f.Write(dataSizeBytes)

	samples16bytes := make([]byte, 0)
	for _, sample := range samples {
		if sample == byte(0x00) {
			samples16bytes = append(samples16bytes, []byte{0x00, 0x80}...)
		} else {
			samples16bytes = append(samples16bytes, []byte{0x00, 0x7f}...)
		}
	}

	f.Write(samples16bytes)

	f.Close()
}
