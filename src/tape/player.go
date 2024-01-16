package tape

import (
	"github.com/gordonklaus/portaudio"
	"io"
)

// Player plays a TZX file as audio samples through the sound card
type Player struct {
	reader   *Reader
	playing  bool
	pause    bool
	stop     bool
	savedPos int64
}

type PlayerInfos struct {
	Playing      bool
	Pause        bool
	CurrentByte  int64
	TotalBytes   int64
	PosPercent   int64
	PosSeconds   int64
	TotalSeconds int64
	FileName     string
	BlockInfo    string
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
			if err = stream.Write(); err != nil {
				//panic(err)
			}
			if err == io.EOF || p.stop {
				break
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
	p.pause = !p.pause
}

func (p *Player) Pause() {
	p.pause = true
}

func (p *Player) Resume() {
	p.pause = false
}

func (p *Player) Stop() {
	p.stop = true
}

func (p *Player) Infos() PlayerInfos {
	return PlayerInfos{
		Playing:      p.playing,
		Pause:        p.pause,
		CurrentByte:  p.reader.Pos(),
		TotalBytes:   p.reader.Size(),
		PosPercent:   p.reader.PosPercent(),
		PosSeconds:   p.reader.PosSeconds(),
		TotalSeconds: p.reader.TotalSeconds(),
		FileName:     p.reader.FileName(),
		BlockInfo:    p.reader.BlockInfo(),
	}
}

func (p *Player) Rewind() {
	_, _ = p.reader.Seek(-50000, 1)
}

func (p *Player) FastForward() {
	_, _ = p.reader.Seek(50000, 1)
}

// SaveCurrentPos saves the current position in memory.
// This position can be seeked later using Player.GoToSavedPos()
// This emulates a real Counter as seen on tape players and can be used to rewind the tape
// at a specific position. Useful for multi-load levels games for example.
func (p *Player) SaveCurrentPos() {
	p.savedPos = p.reader.Pos()
}

// GoToSavedPos set the play position of the tape at the previously stored position
// with Player.SaveCurrentPos. If no position was previously stored, it rewinds the
// tape to beginning.
func (p *Player) GoToSavedPos() error {
	_, err := p.reader.Seek(p.savedPos, 0)
	return err
}
