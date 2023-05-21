package game

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Rock struct {
	X, Y       float32
	Scale      float32
	Dirx, Diry float32
	id         int32
}

const MAX_ROCKS int16 = 100

var Rocks = [MAX_ROCKS]*Rock{}

var original_rock_poly = []float64{

	0, 253,
	89, 121,
	99, 19,
	194, 0,
	343, 76,
	503, 71,
	460, 271,
	510, 410,
	321, 510,
	95, 430,
	0, 253,
}

var rockId int32 = 1

func NewRock(x float32, y float32, scale float32) *Rock {
	dirx := float32(float32(Random_range(-100, 100)) / 300)
	diry := float32(float32(Random_range(-100, 100)) / 300)

	r := Rock{x, y, scale, dirx, diry, rockId}
	rockId++

	log.Printf("rock to add %v\n", r)

	return &r
}

func CentreRockPolyArrayAndFlipCoOrds() {
	for i := 0; i < len(original_rock_poly); i += 2 {
		original_rock_poly[i] = original_rock_poly[i] - float64(RockBounds.Dx())/2
		original_rock_poly[i+1] = float64(RockBounds.Dy()) - original_rock_poly[i+1] - float64(RockBounds.Dy())/2
	}
}

func ResetRocks() {
	for i := 0; i < len(Rocks); i++ {
		Rocks[i] = nil
	}

}

func AddRock(rock *Rock) *Rock {

	before := 0
	for i := 0; i < len(Rocks); i++ {
		if Rocks[i] != nil {
			before++
		}
	}

	if rock.Scale > 10 {
		for i := 0; i < len(Rocks); i++ {
			if Rocks[i] == nil {
				dirx := (rock.Dirx * float32(float32(Random_range(-100, 100))/100) * 20) / 2
				diry := (rock.Diry * float32(float32(Random_range(-100, 100))/100) * 20) / 2

				log.Println(dirx, diry)

				Rocks[i] = NewRock(rock.X+dirx, rock.Y+diry, rock.Scale)
				log.Printf("orig=%v    new=%v \n", rock, Rocks[i])
				break
			}
		}
	} else {
		rock = nil
		log.Println("Rock too small so destroy")
	}
	after := 0
	for i := 0; i < len(Rocks); i++ {
		if Rocks[i] != nil {
			after++
		}
	}
	log.Printf("Rocks before  == %d    ", before)
	log.Printf("Rocks now  == %d\n", after)
	return rock
}

func drawSpaceImage(screen *ebiten.Image, img *ebiten.Image, x, y, width, height float32, rotate int, playerX float32, playerY float32) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(float64(width)/float64(RockBounds.Dx()), float64(height)/float64(RockBounds.Dy()))

	op.GeoM.Translate(-float64(width)/2, -float64(height)/2)

	op.GeoM.Rotate(-float64(rotate%360) * 2 * math.Pi / 360)
	sin := float32(math.Sin(float64(-rotate) * math.Pi / 180.0))
	cos := float32(math.Cos(float64(-rotate) * math.Pi / 180.0))

	x = x - playerX
	y = y - playerY
	coordsX := x*cos + y*sin + PLAYER_SCREEN_X
	coordsY := -x*sin + y*cos + PLAYER_SCREEN_Y

	op.GeoM.Translate(float64(coordsX), float64(ScreenHeight-coordsY))

	screen.DrawImage(img, op)

}

func DrawRocks(screen *ebiten.Image) {

	// draw rocks 4 times to deal with screen wrap around

	drawRocks(screen, 0, ScreenHeight)
	drawRocks(screen, 0, -ScreenHeight)
	drawRocks(screen, ScreenWidth, 0)
	drawRocks(screen, -ScreenWidth, 0)

}

func drawRocks(screen *ebiten.Image, addX float32, addY float32) {
	playerY := XY.Y + addY
	playerX := XY.X + addX

	if playerX < -HalfSpaceWidth {
		playerX = playerX + HalfSpaceWidth*2
	} else if playerX > HalfSpaceWidth {
		playerX = playerX - HalfSpaceWidth*2
	}

	if playerY < -HalfSpaceHeight {
		playerY = playerY + HalfSpaceHeight*2

	} else if playerY > HalfSpaceHeight {
		playerY = playerY - HalfSpaceHeight*2

	}

	for i := 0; i < len(Rocks); i++ {
		rock := Rocks[i]
		if rock != nil {
			rockY := rock.Y + addY
			rockX := rock.X + addX

			if onScreen(rockX, rockY) {

				drawSpaceImage(screen, RockImage, rockX, rockY, rock.Scale, rock.Scale, XY.Rotate, playerX, playerY)

			}
		}
	}
}

func onScreen(x float32, y float32) bool {

	if x >= -ScreenWidth*3 && x <= ScreenWidth*3 &&
		y >= -ScreenHeight*3 && y <= ScreenHeight*3 {
		return true
	}

	return false
}

func UpdateRock() {

	for i := 0; i < len(Rocks); i++ {

		rock := Rocks[i]

		if rock == nil {
			continue
		}
		rock.X = rock.X + rock.Dirx
		rock.Y = rock.Y + rock.Diry

		if rock.X < -HalfSpaceWidth {
			rock.X = rock.X + HalfSpaceWidth*2
		}
		if rock.X > HalfSpaceWidth {
			rock.X = rock.X - HalfSpaceWidth*2
		}
		if rock.Y < -HalfSpaceHeight {
			rock.Y = rock.Y + HalfSpaceHeight*2
		}
		if rock.Y > HalfSpaceHeight {
			rock.Y = rock.Y - HalfSpaceHeight*2
		}

	}

}
