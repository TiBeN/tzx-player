package tape

import (
	"fmt"
	"go.bug.st/serial"
	"log"
	"strconv"
)

// GpioRemoteControl controls the play/pause status of the Player
// according to the Amstrad CPC datacorder relay status by capturing the tension with a GPIO device.
// This code is specific to the Numato Labs 8 channels GPIO module.
// I'm looking in a more vendor-agnostic way to handle this.
type GpioRemoteControl struct {
	port     string
	baudRate int
	ioNb     int
	player   *Player
	start    bool
}

func NewGpioRemoteControl(port string, baudRate int, ioNb int, player *Player) *GpioRemoteControl {
	return &GpioRemoteControl{
		port:     port,
		baudRate: baudRate,
		ioNb:     ioNb,
		player:   player,
	}
}

// Start enable control of the player
func (g *GpioRemoteControl) Start() error {
	g.start = true

	mode := &serial.Mode{
		BaudRate: g.baudRate,
	}
	port, err := serial.Open(g.port, mode)
	if err != nil {
		return err
	}

	command := []byte(fmt.Sprintf("gpio read %d\r", g.ioNb))
	state := 0
	go func() {
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
					g.player.Pause()
				} else {
					g.player.Resume()
				}
			}
			if !g.start {
				if err = port.Close(); err != nil {
					panic(err)
				}
			}
		}
	}()

	return nil
}

// Stop stops the control of the player and free GPIO serial file descriptor
func (g *GpioRemoteControl) Stop() {
	g.start = false
}
