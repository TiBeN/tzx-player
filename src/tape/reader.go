package tape

import (
	"bytes"
	"fmt"
	"github.com/TiBeN/tzx-player/tape/block"
	"math"
)

const TStatePerSecond = 1.0 / 3500000

var AllowedBitDepths []int

func init() {
	AllowedBitDepths = []int{8, 16}
}

// Reader is a TZX tape PCM audio sample io.Reader implementation.
// Its converts block pulse to PCM audio samples
type Reader struct {
	tape         *Tape
	samples      *bytes.Reader
	SamplingRate int
	bitDepth     int
}

func NewReader(tape *Tape, samplingRate int, bitDepth int) (*Reader, error) {
	bitDepthAllowed := false
	for _, b := range AllowedBitDepths {
		if b == bitDepth {
			bitDepthAllowed = true
			break
		}
	}

	if !bitDepthAllowed {
		return nil, fmt.Errorf("unsupported bit depth '%d'", bitDepth)
	}

	r := &Reader{
		tape:         tape,
		SamplingRate: samplingRate,
		bitDepth:     bitDepth,
	}
	r.generateSamples()
	return r, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.samples.Read(p)
}

func (r *Reader) Size() int64 {
	return r.samples.Size()
}

func (r *Reader) Pos() int64 {
	return r.samples.Size() - int64(r.samples.Len())
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	return r.samples.Seek(offset, whence)
}

func (r *Reader) generateSamples() {
	samples := make([]byte, 0)

	for _, b := range r.tape.Blocks {
		samples = append(samples, r.pulsesToSamples(b.Pulses())...)
		samples = append(samples, r.pauseToSamples(b.PauseDuration())...)
	}

	// Prevent tapes without trailing pause to load
	// @TODO: this is a workaround, fix why it does not read up to the end
	samples = append(samples, r.pauseToSamples(2000)...)

	r.samples = bytes.NewReader(samples)
}

func (r *Reader) pulsesToSamples(pulses []block.Pulse) []byte {
	samples := make([]byte, 0)

	for _, pulse := range pulses {
		nbSamples := int(math.Ceil((TStatePerSecond / (1.0 / float64(r.SamplingRate))) * float64(pulse.Length)))
		pulseSamples := make([]byte, 0)
		for i := 0; i < nbSamples; i++ {
			pulseSamples = append(pulseSamples, r.sampleValue(pulse.Level)...)
		}
		samples = append(samples, pulseSamples...)
	}

	return samples
}

func (r *Reader) pauseToSamples(duration int) []byte {
	nbSamples := duration * (r.SamplingRate / 1000)
	samples := make([]byte, 0)
	for i := 0; i < nbSamples; i++ {
		samples = append(samples, r.sampleValue(false)...)
	}
	return samples
}

func (r *Reader) sampleValue(level bool) []byte {
	if r.bitDepth == 8 {
		if !level {
			return []byte{0}
		} else {
			return []byte{255}
		}
	} else { // 16 bit
		if !level {
			return []byte{0x00, 0x80}
		} else {
			return []byte{0x00, 0x7f}
		}
	}
}
