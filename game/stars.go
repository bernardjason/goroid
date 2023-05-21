package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Stars struct {
	screenX, screenY float32
}

var stars = [256]Stars{}

func SetupStars() {

	const randomRange = 32
	i := 0
	for x := -ScreenWidth / 2; x < ScreenWidth/2 && i < len(stars); x += 40 {
		for y := -ScreenHeight / 2; y < ScreenHeight/2 && i < len(stars); y += 40 {

			if Random_range(1, 100) > 90 {

				xx := Random_range(-randomRange, randomRange)
				yy := Random_range(-randomRange, randomRange)
				for xmul := -ScreenWidth * 2; xmul <= ScreenWidth*2; xmul += ScreenWidth {
					for ymul := -ScreenHeight * 2; ymul <= ScreenHeight*2; ymul += ScreenHeight {
						i = add(i, x, xx, xmul, y, yy, ymul)
					}
				}
			}

		}
	}

}

func add(i int, x int, xx int, xmul int, y int, yy int, ymul int) int {
	b := &stars[i]
	b.screenX = float32(x + xx + xmul)
	b.screenY = float32(y + yy + ymul)
	i++
	return i
}

func DrawStars(screen *ebiten.Image, sin float32, cos float32) {
	draw(screen, sin, cos, 0, 0)

}

func draw(screen *ebiten.Image, sin float32, cos float32, addX float32, addY float32) {
	white := color.RGBA{R: 155, G: 155, B: 155, A: 0xff}

	for i := 0; i < len(stars); i++ {
		const width = 1
		const height = 1
		star := &stars[i]

		x1 := star.screenX
		y1 := star.screenY
		x1 = x1 + XY.X
		y1 = y1 + XY.Y
		coordsX := ScreenWidth - (x1*cos + y1*sin + PLAYER_SCREEN_X)
		coordsY := -x1*sin + y1*cos + PLAYER_SCREEN_Y

		if onScreen(coordsX, coordsY) {
			vector.StrokeRect(screen, coordsX, coordsY, width, height, 1, white, true)
		}

	}

}
