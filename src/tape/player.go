package tape

import (
	"github.com/gordonklaus/portaudio"
	"io"
)

type Player struct {
	reader  *Reader
	playing bool
	pause   bool
}

type PlayerInfos struct {
	Playing     bool
	Pause       bool
	CurrentByte int64
	TotalBytes  int64
}

func NewPlayer(reader *Reader) *Player {
	return &Player{
		reader:  reader,
		playing: false,
	}
}

// Start playing the tape
func (p *Player) Start() error {
	if err := portaudio.Initialize(); err != nil {
		return err
	}

	buf := make([]byte, 1000)
	stream, err := portaudio.OpenDefaultStream(
		0,
		1,
		float64(p.reader.SamplingRate),
		len(buf),
		&buf,
	)
	if err != nil {
		return err
	}

	if err = stream.Start(); err != nil {
		return err
	}
	p.playing = true

	// Main playing loop
	go func() {
		for {

			// Handle pause
			if p.pause {
				continue
			}

			_, err = p.reader.Read(buf)
			if err == io.EOF {
				break
			}
			if err = stream.Write(); err != nil {
				//panic(err)
			}
		}

		if err = stream.Stop(); err != nil {
			//panic(err)
		}
		if err = stream.Close(); err != nil {
			//panic(err)
		}

		p.playing = false
	}()

	return nil
}

func (p *Player) TogglePause() {
	// @TODO rework this to be more reactive
	p.pause = !p.pause
}

func (p *Player) Pause() {
	p.pause = true
}

func (p *Player) Resume() {
	p.pause = false
}

func (p *Player) Stop() {

}

func (p *Player) Infos() PlayerInfos {
	return PlayerInfos{
		Playing:     p.playing,
		Pause:       p.pause,
		CurrentByte: p.reader.Pos(),
		TotalBytes:  p.reader.Size(),
	}
	// file name
	// counter [054]
	// position time HH:mm:ss
	// position percent
	// block number
	// block name

	// Infos and time computing made the reader
}

func (p *Player) Rewind() {
	_, _ = p.reader.Seek(-50000, 1)
}

func (p *Player) FastForward() {
	_, _ = p.reader.Seek(50000, 1)
}

// Rewind / Advance

// Counter reset

// Seek to counter 0