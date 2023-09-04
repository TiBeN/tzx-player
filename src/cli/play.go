package cli

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"strconv"
)

type Play struct {
}

func (c *Play) Name() string {
	return "play"
}

func (c *Play) Description() string {
	return "Play a TZX tape"
}

func (c *Play) Usage() string {
	usage := fmt.Sprintf("    Args:\n")
	usage += fmt.Sprintf("      tzx-player play INPUT_TZX_FILE\n")
	usage += fmt.Sprintf("    Options:\n")
	usage += fmt.Sprintf("      %-20sSampling rate (default: %d)\n", "-s int", ConvertDefaultSamplingRate)
	usage += fmt.Sprintf("      %-20sBit depth (default: %d, possibles values: 8 or 16)\n", "-b int", ConvertDefaultBitDepth)
	return usage
}

func (c *Play) Exec(service *tape.Service, args []string) error {
	var tzxFile string

	var err error
	samplingRate := ConvertDefaultSamplingRate
	bitDepth := ConvertDefaultBitDepth

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
		default:
			if tzxFile == "" {
				tzxFile = args[i]
			}
		}
	}

	return service.Play(tzxFile, samplingRate, bitDepth)
}
