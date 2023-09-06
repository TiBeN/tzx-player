package cli

import (
	"errors"
	"fmt"
	"github.com/TiBeN/tzx-player/tape"
	"github.com/eiannone/keyboard"
	"go.bug.st/serial"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
	usage := fmt.Sprintf("    Args:\n")
	usage += fmt.Sprintf("      tzx-player play INPUT_TZX_FILE\n")
	usage += fmt.Sprintf("    Options:\n")
	usage += fmt.Sprintf("      %-20sSampling rate (default: %d)\n", "-s int", ConvertDefaultSamplingRate)
	usage += fmt.Sprintf("      %-20sBit depth (default: %d, possibles values: 8 or 16)\n", "-b int", ConvertDefaultBitDepth)
	usage += fmt.Sprintf("      %-20sSpeed factor: multiply the speed of the tones (experimental) (default: %.1f)\n", "-f float", ConvertDefaultSpeedFactor)
	return usage
}

func (c *Play) Exec(service *tape.Service, args []string) error {
	var tzxFile string

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
			}
		}
	}

	player, err := service.Play(tzxFile, samplingRate, bitDepth, speedFactor)
	if err != nil {
		return err
	}

	// Infos go routine here
	sigs := make(chan os.Signal, 1)

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

	go func() {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()
		for {
			_, key, err := keyboard.GetKey()
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
			if key == keyboard.KeyCtrlC {
				sigs <- syscall.SIGTERM
				break
			}
		}
	}()

	// GPIO remote control handling
	go func() {
		mode := &serial.Mode{
			BaudRate: 19200,
		}
		port, err := serial.Open("/dev/ttyACM0", mode)
		if err != nil {
			panic(err)
		}

		command := []byte("gpio read 1\n\r")
		state := 0
		for {
			_, err = port.Write(command)
			if err != nil {
				panic(err)
			}
			buff := make([]byte, 100)
			for {
				n, err := port.Read(buff)
				if err != nil {
					log.Fatal(err)
				}
				if n < len(buff) {
					break
				}
			}
			val, err := strconv.Atoi(string(buff[len(command)+1 : len(command)+2]))
			if err != nil {
				panic(err)
			}

			if val != state {
				state = val
				if state == 1 {
					player.Pause()
				} else {
					player.Resume()
				}
			}
		}

	}()

	// Lock main thread, handle term signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Print("\n")
	return nil
}
