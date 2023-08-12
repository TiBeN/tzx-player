package main

import "math"

const TStatePerSecond = 1.0 / 3500000

//const TStatePerSecond = (1.0 / 3500000) * (4 / 3.5)

func generatePulseSamples(pulseLength int, samplingRate int, bitDepth int, lowLevel bool) []byte {
	sampleValue := byte(0)
	if !lowLevel {
		sampleValue = byte(255)
	}

	nbSamples := int(math.Ceil((TStatePerSecond / (1.0 / float64(samplingRate))) * float64(pulseLength)))
	samples := make([]byte, nbSamples)
	for i := range samples {
		samples[i] = sampleValue
	}
	return samples
}

func generatePause(pauseDuration int, sampleRate int, bitDepth int) []byte {
	nbSamples := pauseDuration * (sampleRate / 1000) * (bitDepth / 8)
	samples := make([]byte, nbSamples)
	for i := 0; i < nbSamples; i++ {
		samples[i] = 0
	}
	return samples
}
