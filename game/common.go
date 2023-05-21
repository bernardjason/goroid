package game

import (
	"bytes"
	"embed"
	"image"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	EnemyImage, RedImage, RockImage, ShipImage *ebiten.Image
	EnemyBounds, RockBounds, ShipBounds        image.Rectangle
)

const (
	ScreenWidth     = 320
	ScreenHeight    = 320
	HalfSpaceWidth  = ScreenWidth * 1
	HalfSpaceHeight = ScreenHeight * 1
)

const PLAYER_SCREEN_Y = ScreenHeight / 2
const PLAYER_SCREEN_X = ScreenWidth / 2

//go:embed texture.png
var textureFile embed.FS

//go:embed pixelship.png
var pixelFile embed.FS

//go:embed red.png
var redRock embed.FS

//go:embed enemy.png
var enemyImageFile embed.FS

func LoadImages() {
	file, err := textureFile.ReadFile("texture.png")
	if err != nil {
		log.Fatal(err)
	}

	RockImage = loadImage(&file)
	RockBounds = RockImage.Bounds()

	redfile, err := redRock.ReadFile("red.png")
	if err != nil {
		log.Fatal(err)
	}

	RedImage = loadImage(&redfile)

	pixelship, err := pixelFile.ReadFile("pixelship.png")
	if err != nil {
		log.Fatal(err)
	}
	ShipImage = loadImage(&pixelship)
	ShipBounds = ShipImage.Bounds()

	enemyShip, err := enemyImageFile.ReadFile("enemy.png")
	if err != nil {
		log.Fatal(err)
	}
	EnemyImage = loadImage(&enemyShip)
	EnemyBounds = EnemyImage.Bounds()
}

func loadImage(imageFileName *[]byte) *ebiten.Image {

	img, _, err := image.Decode(bytes.NewReader(*imageFileName))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func Random_range(min int, max int) int {

	return rand.Intn(max-min) + min
}
func CheckRockHit(sound *Player) bool {
	hitWall := false
	hitARock := false
	player_ship_poly := CentreShipAbsolute()
	bulletHitPlayer(player_ship_poly)

hit:
	for i := 0; i < len(Rocks); i++ {

		rock := Rocks[i]

		if rock == nil {
			continue
		}

		hitWall, Rocks[i], hitARock = playerCollisionLogicForARock(rock, player_ship_poly)

		if hitWall {
			Rocks[i] = nil
			sound.playExplode()
			Explode(XY.X, XY.Y)
			break hit

		}
		if hitARock {
			sound.playExplode()
		}

	}

	return hitWall
}

func playerCollisionLogicForARock(rock *Rock, player_ship []float64) (bool, *Rock, bool) {
	hitWall := false
	var poly []float64 = make([]float64, len(original_rock_poly))
	copy(poly, original_rock_poly)
	scale := float64(rock.Scale) / float64(RockBounds.Dx())

	for i := 0; i < len(poly); i += 2 {

		poly[i] = poly[i]*scale + float64(rock.X)
		poly[i+1] = poly[i+1]*scale + float64(rock.Y)
	}

	hitWall = playerHitRock(rock, player_ship, poly)
	addedRock, hitARock := rockHitByPlayer(rock, poly)
	return hitWall, addedRock, hitARock
}

func playerHitRock(rock *Rock, player_ship []float64, poly []float64) bool {

	hitPlayer := false

	for i := 0; i < len(player_ship); i = i + 2 {
		x1 := float32(player_ship[i])
		y1 := float32(player_ship[i+1])

		checked := pointInside(float64(x1), float64(y1), poly)

		if checked {
			log.Printf("*** HIT ROCK %d *** %0.2f,%0.2f rotation %d \n", rock.id, x1, y1, XY.Rotate) //, away)
			hitPlayer = true
			XY.Lives = XY.Lives - 1

			break
		}
	}
	return hitPlayer
}

func EnemyHitbyBullet() {
	for ei := 0; ei < len(Enemies); ei++ {
		enemy := Enemies[ei]
		if enemy != nil {
			poly := CentreEnemyAbsolute(enemy.X, enemy.Y)
			for bi := 0; bi < len(bullets); bi++ {
				b := &bullets[bi]
				if b.slotOccupied && b.players {
					checked := pointInside(float64(b.absX), float64(b.absY), *poly)
					if checked {
						b.slotOccupied = false
						XY.Score = XY.Score + 100
						log.Printf("EnemyHitbyBullet -Hit enemy\n")
						Enemies[ei] = nil
						Explode(b.screenX, b.screenY)
					}
				}
			}
		}
	}
}

func bulletHitPlayer(poly []float64) {

	for bi := 0; bi < len(bullets); bi++ {
		b := &bullets[bi]
		if b.slotOccupied && !b.players {

			checked := pointInside(float64(b.absX), float64(b.absY), poly)
			if checked {
				b.slotOccupied = false
				XY.Lives = XY.Lives - 1
				log.Printf("bulletHitPlayer - Player hit !!!!!!!!\n")

				Explode(b.screenX, b.screenY)
			}
		}
	}
}

func rockHitByPlayer(rock *Rock, poly []float64) (*Rock, bool) {
	var addedRock *Rock = rock
	hit := false
	for bi := 0; bi < len(bullets); bi++ {
		b := &bullets[bi]
		if b.slotOccupied {
			checked := pointInside(float64(b.absX), float64(b.absY), poly)
			if checked {
				b.slotOccupied = false
				XY.Score = XY.Score + 1
				log.Printf("rockHitByPlayer - Hit rock\n")
				hit = true
				rock.Scale = rock.Scale / 2
				addedRock = AddRock(rock)
				Explode(b.screenX, b.screenY)
				Explode(rock.X, rock.Y)
			}
		}
	}
	return addedRock, hit
}

// http://www.ariel.com.au/a/python-point-int-poly.html
func pointInside(x float64, y float64, poly []float64) bool {

	n := len(poly)
	inside := false

	p1x := poly[0]
	p1y := poly[1]
	xinters := 0.0
	for i := 0; i <= n; i += 2 {
		p2x := poly[i%n]
		p2y := poly[(i%n)+1]
		if y > math.Min(p1y, p2y) {
			if y <= math.Max(p1y, p2y) {
				if x <= math.Max(p1x, p2x) {
					if p1y != p2y {
						xinters = (y-p1y)*(p2x-p1x)/(p2y-p1y) + p1x
					}

					if p1x == p2x || x <= xinters {
						inside = !inside
					}
				}
			}
		}
		p1x, p1y = p2x, p2y
	}
	return inside
}

func CalculateDistanceBetweenPointsWithoutSqrt(
	x1 float32,
	y1 float32,
	x2 float32,
	y2 float32) float32 {

	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)

}
