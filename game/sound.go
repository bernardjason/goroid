package game

import (
	"bytes"
	"embed"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 48000
)

type Player struct {
	audioContext *audio.Context
	//audioPlayer  *audio.Player
	seBytes []byte
	seCh    chan []byte

	explodeChan   chan []byte
	explodeByte   []byte
	explodePlayer *audio.Player
	volume128     int
}

//go:embed cannon.wav
var cannon embed.FS

//go:embed explosion.wav
var explode embed.FS

func InitialiseSound() *Player {
	p := &Player{
		audioContext: audio.NewContext(sampleRate),
		seCh:         make(chan []byte),
		explodeChan:  make(chan []byte),
		volume128:    128,
	}

	go func() {
		file, err := cannon.ReadFile("cannon.wav")

		if err != nil {
			log.Fatal(err)
		}

		s, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(file))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := io.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		p.seCh <- b
	}()
	go func() {
		file, err := explode.ReadFile("explosion.wav")

		if err != nil {
			log.Fatal(err)
		}

		s, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(file))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := io.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		p.explodeChan <- b
	}()

	return p
}

func (p *Player) playShot() {
	if p.seBytes == nil {
		// Bytes for the SE is not loaded yet.
		return
	}

	sePlayer := p.audioContext.NewPlayerFromBytes(p.seBytes)
	sePlayer.Play()
}
func (p *Player) playExplode() {

	if p.explodeByte == nil {
		// Bytes for the SE is not loaded yet.
		return
	}

	if p.explodePlayer == nil || !p.explodePlayer.IsPlaying() {
		sePlayer := p.audioContext.NewPlayerFromBytes(p.explodeByte)
		p.explodePlayer = sePlayer
		sePlayer.Play()
	}
}

func (p *Player) SoundReadyToPlay() {

	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}
	select {
	case p.explodeByte = <-p.explodeChan:
		close(p.explodeChan)
		p.explodeChan = nil
	default:
	}

}
