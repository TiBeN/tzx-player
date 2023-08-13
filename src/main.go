package main

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"os"
	"time"
)

func main() {
	params, err := parseArgs()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		usage()
		os.Exit(3)
	}

	t, err := tape.NewTape(params.TzxTapeInputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error()+"\n")
		usage()
		os.Exit(4)
	}

	for _, infoLine := range t.Info() {
		fmt.Println(infoLine)
	}

	start := time.Now()
	samples := t.Samples(44100, 8)
	end := time.Now()
	fmt.Printf("Generated %d samples in %s", len(samples), end.Sub(start))

	writeToWavFile(params.OutputFile, 44100, 8, samples)
}

func parseArgs() (*Parameters, error) {
	args := os.Args
	params := Parameters{}

	if len(args) != 3 {
		return nil, errors.New("no TZX tape input file specified")
	}

	params.TzxTapeInputFile = args[1]
	params.OutputFile = args[2]
	return &params, nil
}

func usage() {
	fmt.Println("")
	fmt.Println("Usage: tzx-player <input-file>.(tzx|cdt) <output-file>.wav")
}
