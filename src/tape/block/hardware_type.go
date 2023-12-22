package block

import (
	"fmt"
	"os"
)

type HardwareTypeIds struct {
	name        string
	hardwareIds map[byte]string
}

var HardwareList map[byte]HardwareTypeIds

var HardwareInfosDescription map[byte]string

func init() {
	HardwareList = map[byte]HardwareTypeIds{
		0x00: {
			name: "Computers",
			hardwareIds: map[byte]string{
				0x00: "ZX Spectrum 16k",
				0x01: "ZX Spectrum 48k, Plus",
				0x02: "ZX Spectrum 48k ISSUE 1",
				0x03: "ZX Spectrum 128k +(Sinclair)",
				0x04: "ZX Spectrum 128k +2 (grey case)",
				0x05: "ZX Spectrum 128k +2A, +3",
				0x06: "Timex Sinclair TC-2048",
				0x07: "Timex Sinclair TS-2068",
				0x08: "Pentagon 128",
				0x09: "Sam Coupe",
				0x0A: "Didaktik M",
				0x0B: "Didaktik Gama",
				0x0C: "ZX-80",
				0x0D: "ZX-81",
				0x0E: "ZX Spectrum 128k, Spanish version",
				0x0F: "ZX Spectrum, Arabic version",
				0x10: "Microdigital TK 90-X",
				0x11: "Microdigital TK 95",
				0x12: "Byte",
				0x13: "Elwro 800-3",
				0x14: "ZS Scorpion 256",
				0x15: "Amstrad CPC 464",
				0x16: "Amstrad CPC 664",
				0x17: "Amstrad CPC 6128",
				0x18: "Amstrad CPC 464+",
				0x19: "Amstrad CPC 6128+",
				0x1A: "Jupiter ACE",
				0x1B: "Enterprise",
				0x1C: "Commodore 64",
				0x1D: "Commodore 128",
				0x1E: "Inves Spectrum+",
				0x1F: "Profi",
				0x20: "GrandRomMax",
				0x21: "Kay 1024",
				0x22: "Ice Felix HC 91",
				0x23: "Ice Felix HC 2000",
				0x24: "Amaterske RADIO Mistrum",
				0x25: "Quorum 128",
				0x26: "MicroART ATM",
				0x27: "MicroART ATM Turbo 2",
				0x28: "Chrome",
				0x29: "ZX Badaloc",
				0x2A: "TS-1500",
				0x2B: "Lambda",
				0x2C: "TK-65",
				0x2D: "ZX-97",
			},
		},
		0x01: {
			name: "External storage",
			hardwareIds: map[byte]string{
				0x00: "ZX Microdrive",
				0x01: "Opus Discovery",
				0x02: "MGT Disciple",
				0x03: "MGT Plus-D",
				0x04: "Rotronics Wafadrive",
				0x05: "TR-DOS (BetaDisk)",
				0x06: "Byte Drive",
				0x07: "Watsford",
				0x08: "FIZ",
				0x09: "Radofin",
				0x0A: "Didaktik disk drives",
				0x0B: "BS-DOS (MB-02)",
				0x0C: "ZX Spectrum +3 disk drive",
				0x0D: "JLO (Oliger) disk interface",
				0x0E: "Timex FDD3000",
				0x0F: "Zebra disk drive",
				0x10: "Ramex Millenia",
				0x11: "Larken",
				0x12: "Kempston disk interface",
				0x13: "Sandy",
				0x14: "ZX Spectrum +3e hard disk",
				0x15: "ZXATASP",
				0x16: "DivIDE",
				0x17: "ZXCF",
			},
		},
		0x02: {
			name: "ROM/RAM type add-ons",
			hardwareIds: map[byte]string{
				0x00: "Sam Ram",
				0x01: "Multiface ONE",
				0x02: "Multiface 128k",
				0x03: "Multiface +3",
				0x04: "MultiPrint",
				0x05: "MB-02 ROM/RAM expansion",
				0x06: "SoftROM",
				0x07: "1k",
				0x08: "16k",
				0x09: "48k",
				0x0A: "Memory in 8-16k used",
			},
		},
		0x03: {
			name: "Sound devices",
			hardwareIds: map[byte]string{
				0x00: "Classic AY hardware (compatible with 128k ZXs)",
				0x01: "Fuller Box AY sound hardware",
				0x02: "Currah microSpeech",
				0x03: "SpecDrum",
				0x04: "AY ACB stereo (A+C=left, B+C=right); Melodik",
				0x05: "AY ABC stereo (A+B=left, B+C=right)",
				0x06: "RAM Music Machine",
				0x07: "Covox",
				0x08: "General Sound",
				0x09: "Intec Electronics Digital Interface B8001",
				0x0A: "Zon-X AY",
				0x0B: "QuickSilva AY",
				0x0C: "Jupiter ACE",
			},
		},
		0x04: {
			name: "Joysticks",
			hardwareIds: map[byte]string{
				0x00: "Kempston",
				0x01: "Cursor, Protek, AGF",
				0x02: "Sinclair 2 Left (12345)",
				0x03: "Sinclair 1 Right (67890)",
				0x04: "Fuller",
			},
		},
		0x05: {
			name: "Mice",
			hardwareIds: map[byte]string{
				0x00: "AMX mouse",
				0x01: "Kempston mouse",
			},
		},
		0x06: {
			name: "Other controllers",
			hardwareIds: map[byte]string{
				0x00: "Trickstick",
				0x01: "ZX Light Gun",
				0x02: "Zebra Graphics Tablet",
				0x03: "Defender Light Gun",
			},
		},
		0x07: {
			name: "Serial ports",
			hardwareIds: map[byte]string{
				0x00: "ZX Interface 1",
				0x01: "ZX Spectrum 128k",
			},
		},
		0x08: {
			name: "Parallel ports",
			hardwareIds: map[byte]string{
				0x00: "Kempston S",
				0x01: "Kempston E",
				0x02: "ZX Spectrum +3",
				0x03: "Tasman",
				0x04: "DK'Tronics",
				0x05: "Hilderbay",
				0x06: "INES Printerface",
				0x07: "ZX LPrint Interface 3",
				0x08: "MultiPrint",
				0x09: "Opus Discovery",
				0x0A: "Standard 8255 chip with ports 31,63,95",
			},
		},
		0x09: {
			name: "Printers",
			hardwareIds: map[byte]string{
				0x00: "ZX Printer, Alphacom 32 & compatibles",
				0x01: "Generic printer",
				0x02: "EPSON compatible",
			},
		},
		0x0A: {
			name: "Modems",
			hardwareIds: map[byte]string{
				0x00: "Prism VTX 5000",
				0x01: "T/S 2050 or Westridge 2050",
			},
		},
		0x0B: {
			name: "Digitizers",
			hardwareIds: map[byte]string{
				0x00: "RD Digital Tracer",
				0x01: "DK'Tronics Light Pen",
				0x02: "British MicroGraph Pad",
				0x03: "Romantic Robot Videoface",
			},
		},
		0x0C: {
			name: "Network adapters",
			hardwareIds: map[byte]string{
				0x00: "ZX Interface 1",
			},
		},
		0x0D: {
			name: "Keyboards & keypads",
			hardwareIds: map[byte]string{
				0x00: "Keypad for ZX Spectrum 128k",
			},
		},
		0x0E: {
			name: "AD/DA converters",
			hardwareIds: map[byte]string{
				0x00: "Harley Systems ADC 8.2",
				0x01: "Blackboard Electronics",
			},
		},
		0x0F: {
			name: "EPROM programmers",
			hardwareIds: map[byte]string{
				0x00: "Orme Electronics",
			},
		},
		0x10: {
			name: "Graphics",
			hardwareIds: map[byte]string{
				0x00: "WRX Hi-Res",
				0x01: "G007",
				0x02: "Memotech",
				0x03: "Lambda Colour",
			},
		},
	}

	HardwareInfosDescription = map[byte]string{
		0x00: "The tape RUNS on this machine or with this hardware, but may or may not use the hardware or special features of the machine.",
		0x01: "The tape USES the hardware or special features of the machine, such as extra memory or a sound chip.",
		0x02: "The tape RUNS but it DOESN'T use the hardware or special features of the machine.",
		0x03: "The tape DOESN'T RUN on this machine or with this hardware.",
	}
}

type HardwareType struct {
	hardwareInfos []HardwareInfo
}

type HardwareInfo struct {
	hwType byte
	id     byte
	info   byte
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

	for i := 0; i < int(machineNb[0]); i++ {
		hwInfo := make([]byte, 3)
		if _, err := tzxFile.Read(hwInfo); err != nil {
			return err
		}
		h.hardwareInfos = append(h.hardwareInfos, HardwareInfo{
			hwType: hwInfo[0],
			id:     hwInfo[1],
			info:   hwInfo[2],
		})
	}

	return nil
}

func (h *HardwareType) Info() [][]string {
	info := make([][]string, 0)
	for _, t := range h.hardwareInfos {
		info = append(info, []string{
			HardwareList[t.hwType].name,
			fmt.Sprintf(
				"%s (%s)",
				HardwareList[t.hwType].hardwareIds[t.id],
				HardwareList[t.hwType].hardwareIds[t.info],
			)})
	}
	return info
}

func (h *HardwareType) Pulses() []Pulse {
	return make([]Pulse, 0)
}

func (h *HardwareType) PauseDuration() int {
	return 0
}
