package tape

// Service is the main interface for actions related to a TZX Tape
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ConvertToWavFile converts the given TZX tape file into an audio PCM wav file
func (s *Service) ConvertToWavFile(tzxFile string, outputFile string, samplingRate int, bitDepth int) error {
	tape, err := NewTape(tzxFile)
	if err != nil {
		return err
	}

	writeToWavFile(outputFile, 44100, 8, tape.Samples(44100, 8))
	return nil
}
