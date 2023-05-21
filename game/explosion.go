package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Explosion struct {
	screenX, screenY float32
	dirX, dirY       float32
	ttl              int16
	notFree          bool
}

const PARTICLE_TTL int16 = 30

var particles = [256]Explosion{}

func UpdateParticles() {
	for i := len(particles) - 1; i >= 0; i-- {
		b := &particles[i]
		if b.notFree {
			b.screenX = b.screenX + b.dirX
			b.screenY = b.screenY + b.dirY

			b.ttl = b.ttl - 1
			if b.ttl < 0 {
				b.notFree = false
			}
		}
	}
}

func Explode(x float32, y float32) {
	addParticles := 20
	for i := 0; i < len(particles) && addParticles > 0; i++ {
		b := &particles[i]
		if !b.notFree {

			addX := float32(Random_range(-100, 100)) / 200
			addY := float32(Random_range(-100, 100)) / 200

			b.dirX = float32(addX)
			b.dirY = float32(addY)
			b.screenX = x + b.dirX
			b.screenY = y + b.dirY
			b.ttl = PARTICLE_TTL
			b.notFree = true
			addParticles--
		}
	}
}

func DrawExplosions(screen *ebiten.Image, sin float32, cos float32) {
	colours := [...]color.Color{color.RGBA{R: 255, G: 255, B: 255, A: 0xff},
		color.RGBA{R: 255, G: 255, B: 0, A: 0xff},
		color.RGBA{R: 255, G: 0, B: 0, A: 0xff}}
	const width = 1
	const height = 1

	colourRange := len(colours)
	colour := Random_range(1, 10)
	for i := 0; i < len(particles); i++ {

		b := &particles[i]
		if b.notFree {

			x1 := b.screenX
			y1 := b.screenY
			x1 = x1 - XY.X
			y1 = y1 - XY.Y
			coordsX := x1*cos + y1*sin + PLAYER_SCREEN_X
			coordsY := -x1*sin + y1*cos + PLAYER_SCREEN_Y
			colour++

			vector.DrawFilledRect(screen, coordsX, ScreenHeight-coordsY, width*4, height*4, colours[colour%colourRange], true)

		}

	}
}
