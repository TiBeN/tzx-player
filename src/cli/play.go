package cli

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"github.com/eiannone/keyboard"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"
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
	usage := fmt.Sprintln("    Args:")
	usage += fmt.Sprintln("      tzx-player play INPUT_TZX_FILE")
	usage += fmt.Sprintln("    Options:")
	usage += fmt.Sprintf("      %-20sSampling rate (default: %d)\n", "-s int", ConvertDefaultSamplingRate)
	usage += fmt.Sprintf("      %-20sBit depth (default: %d, possibles values: 8 or 16)\n", "-b int", ConvertDefaultBitDepth)
	usage += fmt.Sprintf("      %-20sEnable tape remote control using a GPIO device. Support only Numato labs GPIO Modules for now. Exemple: -g /dev/ttyACM0:9600:1\n", "-g port:baud:ionb")
	usage += fmt.Sprintf("      %-20sSpeed factor: multiply the speed of the tones (experimental) (default: %.1f)\n", "-f float", ConvertDefaultSpeedFactor)
	usage += fmt.Sprintln("   Player control keystrokes:")
	usage += fmt.Sprintln("       Space : Toggle play/pause")
	usage += fmt.Sprintln("       Right arrow : Fast forward")
	usage += fmt.Sprintln("       Left arrow : Rewind")
	usage += fmt.Sprintln("       g : Set tape to last saved position")
	usage += fmt.Sprintln("       p : Pause")
	usage += fmt.Sprintln("       s : Save current tape position")
	usage += fmt.Sprintln("       g : Set tape to last saved position")

	return usage
}

func (c *Play) Exec(service *tape.Service, args []string) error {
	var tzxFile string

	var err error
	samplingRate := ConvertDefaultSamplingRate
	bitDepth := ConvertDefaultBitDepth
	speedFactor := ConvertDefaultSpeedFactor
	enableGpio := false
	gpioPort := ""
	gpioBaudRate := 0
	gpioNb := 0

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
		case "-g":
			enableGpio = true
			if i == len(args)-1 {
				return fmt.Errorf("missing -g argument")
			}
			argRegex := regexp.MustCompile("([^:]+):([^:]+):([^:]+)")
			values := argRegex.FindStringSubmatch(args[i+1])
			if values == nil {
				return errors.New("-g argument is not valid")
			}
			gpioPort = values[1]
			gpioBaudRate, err = strconv.Atoi(values[2])
			if err != nil {
				return errors.New("-g argument: not a valid baud rate value")
			}
			gpioNb, err = strconv.Atoi(values[3])
			if err != nil {
				return errors.New("-g argument: not a valid IO port number value")
			}
			i++
		default:
			if tzxFile == "" {
				tzxFile = args[i]
			}
		}
	}

	player, err := service.Play(tzxFile, samplingRate, bitDepth, speedFactor)
	if err != nil {
		return err
	}
	defer player.Stop()

	sigs := make(chan os.Signal, 1)

	// Infos status bar
	go func() {
		infosTicker := time.NewTicker(time.Duration(60) * time.Millisecond)
		for {
			<-infosTicker.C
			playerInfos := player.Infos()
			if !playerInfos.Playing {
				sigs <- syscall.SIGTERM
			}

			playStatus := "\u23F5"
			if playerInfos.Pause {
				playStatus = "\u23F8"
			}

			totalTime := time.Unix(playerInfos.TotalSeconds, 0)
			currentTime := time.Unix(playerInfos.PosSeconds, 0)

			fmt.Printf(
				"\r\033[K%s %s - %s / %s (%d%%) - Block: %s",
				playStatus,
				filepath.Base(playerInfos.FileName),
				currentTime.UTC().Format("15:04:05"),
				totalTime.UTC().Format("15:04:05"),
				playerInfos.PosPercent,
				playerInfos.BlockInfo,
			)
		}
	}()

	// Handle keyboard shortcuts
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeySpace {
				player.TogglePause()
			}
			if key == keyboard.KeyArrowLeft {
				player.Rewind()
			}
			if key == keyboard.KeyArrowRight {
				player.FastForward()
			}
			if char == 's' {
				player.SaveCurrentPos()
			}
			if char == 'g' {
				if err := player.GoToSavedPos(); err != nil {
					panic(err)
				}
			}
			if key == keyboard.KeyCtrlC {
				sigs <- syscall.SIGTERM
				break
			}
		}
	}()

	// GPIO remote control handling
	if enableGpio {
		gpio := tape.NewGpioRemoteControl(gpioPort, gpioBaudRate, gpioNb, player)
		if err = gpio.Start(); err != nil {
			return err
		}
		defer gpio.Stop()
	}

	// Lock main thread, handle term signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Print("\n")
	return nil
}
