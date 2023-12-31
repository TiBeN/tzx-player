package cli

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"strconv"
)

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
	usage += fmt.Sprintf("      %-20sSpeed factor: multiply the speed of the tones (experimental) (default: %.1f)\n", "-f float", ConvertDefaultSpeedFactor)
	return usage
}

func (c *Convert) Exec(service *tape.Service, args []string) error {
	var tzxFile string
	var outputFile string
	var err error
	samplingRate := ConvertDefaultSamplingRate
	bitDepth := ConvertDefaultBitDepth
	speedFactor := ConvertDefaultSpeedFactor

	// Parse args
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s":
			if i == len(args)-1 {
				return errors.New("missing -s argument")
			}
			samplingRate, err = strconv.Atoi(args[i+1])
			if err != nil {
				return errors.New("-s argument is not a valid number")
			}
			i++
		case "-b":
			if i == len(args)-1 {
				return fmt.Errorf("missing -b argument")
			}
			bitDepth, err = strconv.Atoi(args[i+1])
			if err != nil {
				return errors.New("-s argument is not a valid number")
			}
			i++
		case "-f":
			if i == len(args)-1 {
				return fmt.Errorf("missing -f argument")
			}
			speedFactor, err = strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				return errors.New("-f argument is not a valid number")
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

	generationTime, err := service.ConvertToWavFile(tzxFile, outputFile, samplingRate, bitDepth, speedFactor)

	if err == nil {
		fmt.Printf("Generation time: %s\n", generationTime)
	}

	return err
}
