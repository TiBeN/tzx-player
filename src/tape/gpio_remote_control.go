package tape

import (
	"fmt"
	"go.bug.st/serial"
	"log"
	"strconv"
)

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

func (g *GpioRemoteControl) Stop() {
	g.start = false
}
