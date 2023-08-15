package cli

import (
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"strconv"
)

const ConvertDefaultSamplingRate = 44100
const ConvertDefaultBitDepth = 8

type Convert struct {
}

func (c *Convert) Name() string {
	return "convert"
}

func (c *Convert) Description() string {
	return "Convert TZX tape to an audio PCM Wav file"
}

func (c *Convert) Usage() string {
	usage := fmt.Sprintf("    Args:\n")
	usage += fmt.Sprintf("      tzx-player convert INPUT_TZX_FILE OUTPUT_WAV_FILE\n")
	usage += fmt.Sprintf("    Options:\n")
	usage += fmt.Sprintf("      %-20sSampling rate (default: %d)\n", "-s int", ConvertDefaultSamplingRate)
	usage += fmt.Sprintf("      %-20sBit depth (default: %d, possibles values: 8 or 16)\n", "-b int", ConvertDefaultBitDepth)
	return usage
}

func (c *Convert) Exec(service *tape.Service, args []string) error {
	var tzxFile string
	var outputFile string
	var err error
	samplingRate := ConvertDefaultSamplingRate
	bitDepth := ConvertDefaultBitDepth

	// Parse args
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s":
			if i == len(args)-1 {
				return fmt.Errorf("missing -s argument")
			}
			samplingRate, err = strconv.Atoi(args[i+1])
			if err != nil {
				return fmt.Errorf("-s argument is not a valid number")
			}
			i++
		case "-b":
			if i == len(args)-1 {
				return fmt.Errorf("missing -b argument")
			}
			bitDepth, err = strconv.Atoi(args[i+1])
			if err != nil || (bitDepth != 8 && bitDepth != 16) {
				return fmt.Errorf("-s argument is not a valid number (8 or 16 supported)")
			}
			i++
		default:
			if tzxFile == "" {
				tzxFile = args[i]
			} else {
				outputFile = args[i]
			}
		}
	}

	return service.ConvertToWavFile(tzxFile, outputFile, samplingRate, bitDepth)
}
