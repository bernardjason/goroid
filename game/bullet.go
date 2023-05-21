package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Bullet struct {
	slotOccupied     bool
	screenX, screenY float32
	absX, absY       float32
	dirx, diry       float32
	ttl              int32
	players          bool
}

const PLAYER_TTL int32 = 120
const ENEMY_TTL int32 = 240

var bullets = [100]Bullet{}

func UpdateBullets() {
	for i := len(bullets) - 1; i >= 0; i-- {
		b := &bullets[i]
		if b.slotOccupied {
			var scale float32 = 1.0
			if !b.players {
				scale = 0.35
			}
			b.screenX = b.screenX + b.dirx*scale
			b.screenY = b.screenY + b.diry*scale
			b.absX = b.screenX
			b.absY = b.screenY

			if b.absX < -HalfSpaceWidth {
				b.absX = b.absX + HalfSpaceWidth*2
			} else if b.absX > HalfSpaceWidth {
				b.absX = b.absX - HalfSpaceWidth*2
			}

			if b.absY < -HalfSpaceHeight {
				b.absY = b.absY + HalfSpaceHeight*2

			} else if b.absY > HalfSpaceHeight {
				b.absY = b.absY - HalfSpaceHeight*2

			}

			b.ttl = b.ttl - 1
			if b.ttl < 0 {
				b.slotOccupied = false

			}
		}
	}
}

func ResetBullets() {
	for i := len(bullets) - 1; i >= 0; i-- {
		b := &bullets[i]
		b.slotOccupied = false
	}
}

func FireBullet(height float32, x float32, y float32, rotate float32, playerBullet bool) {
	sin := math.Sin(float64(rotate) * math.Pi / 180.0)
	cos := math.Cos(float64(rotate) * math.Pi / 180.0)
foundFreeSlot:
	for i := 0; i < len(bullets); i++ {
		b := &bullets[i]
		if !b.slotOccupied {

			x1 := 0.0
			y1 := float64(height)

			coordsX := x1*cos + y1*sin
			coordsY := -x1*sin + y1*cos

			b.screenX = x + float32(coordsX*0.75)
			b.screenY = y + float32(coordsY*0.75)
			b.absX = b.screenX
			b.absY = b.screenY
			b.dirx = float32(sin * 3)
			b.diry = float32(cos * 3)
			if playerBullet {
				b.ttl = PLAYER_TTL
			} else {
				b.ttl = ENEMY_TTL
			}
			b.slotOccupied = true
			b.players = playerBullet
			break foundFreeSlot
		}
	}
}

func DrawBullets(screen *ebiten.Image, sin float32, cos float32) {
	const width = 2
	const height = 2

	enemyColour := color.RGBA{R: 0, G: 255, B: 0, A: 0xff}
	playerColour := color.RGBA{R: 255, G: 155, B: 0, A: 0xff}

	for i := 0; i < len(bullets); i++ {

		b := &bullets[i]
		if b.slotOccupied {
			x1 := b.screenX
			y1 := b.screenY
			x1 = x1 - XY.X
			y1 = y1 - XY.Y
			coordsX := x1*cos + y1*sin + PLAYER_SCREEN_X
			coordsY := -x1*sin + y1*cos + PLAYER_SCREEN_Y
			if coordsX >= -width && coordsX <= ScreenWidth+width {
				if b.players {
					vector.DrawFilledRect(screen, coordsX, ScreenHeight-coordsY, width, height, playerColour, true)
				} else {
					vector.DrawFilledRect(screen, coordsX, ScreenHeight-coordsY, width, height, enemyColour, true)
				}
			}
		}

	}
}
