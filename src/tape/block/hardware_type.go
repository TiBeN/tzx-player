package block

import "os"

type HardwareType struct {
}

func (h *HardwareType) Id() byte {
	return 0x33
}

func (h *HardwareType) Name() string {
	return "Hardware Type"
}

func (h *HardwareType) Read(tzxFile *os.File) error {
	machineNb := make([]byte, 1)
	if _, err := tzxFile.Read(machineNb); err != nil {
		return err
	}

	hwInfo := make([]byte, machineNb[0]*3)
	if _, err := tzxFile.Read(hwInfo); err != nil {
		return err
	}

	return nil
}

func (h *HardwareType) Info() string {
	return "not implemented"
}

func (h *HardwareType) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (h *HardwareType) PauseDuration() int {
	return 0
}
