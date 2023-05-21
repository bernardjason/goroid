package game

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	X, Y       float32
	Dirx, Diry float32
	moveAway   bool
}

const MAX_ENEMIY int16 = 3

var Enemies = [MAX_ENEMIY]*Enemy{}

var enemy_poly = []float64{
	0, 7,
	6, 0,
	25, 0,
	31, 7,
	25, 15,
	15, 11,
	6, 15,
	0, 7,
}

func CentreEnemyAbsolute(x float32, y float32) *[]float64 {
	var poly []float64 = make([]float64, len(enemy_poly))
	copy(poly, enemy_poly)
	for i := 0; i < len(poly); i += 2 {
		poly[i] = float64(x) + poly[i] - float64(EnemyBounds.Dx())/2
		poly[i+1] = float64(y) + float64(EnemyBounds.Dy()) - poly[i+1] - float64(EnemyBounds.Dy())/2
	}
	return &poly
}

func NewEnemy(x float32, y float32) *Enemy {
	dirx, diry := setDirectionToPlayer(x, y)

	r := Enemy{x, y, dirx, diry, false}

	log.Printf("enemy to add %v\n", r)

	return &r
}

func setDirectionToPlayer(x float32, y float32) (float32, float32) {
	var dirx float32 = 0
	var diry float32 = 0
	if x > XY.X {
		dirx = -1
	} else if x < XY.X {
		dirx = 1
	}

	if y > XY.Y {
		diry = -1
	} else if y < XY.Y {
		diry = 1
	}
	return dirx, diry
}

func AddEnemy(x float32, y float32) *Enemy {

	var enemy *Enemy = nil

	for i := 0; i < len(Enemies); i++ {
		if Enemies[i] == nil {
			Enemies[i] = NewEnemy(x, y)
			break
		}
	}

	return enemy
}

func drawEnemyImage(screen *ebiten.Image, img *ebiten.Image, x, y, width, height float32, rotate int, playerX float32, playerY float32) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(width)/2, -float64(height)/2)

	sin := float32(math.Sin(float64(-rotate) * math.Pi / 180.0))
	cos := float32(math.Cos(float64(-rotate) * math.Pi / 180.0))

	x = x - playerX
	y = y - playerY
	coordsX := x*cos + y*sin + PLAYER_SCREEN_X
	coordsY := -x*sin + y*cos + PLAYER_SCREEN_Y

	op.GeoM.Translate(float64(coordsX), float64(ScreenHeight-coordsY))

	screen.DrawImage(img, op)

}

func DrawEnemy(screen *ebiten.Image, sin float32, cos float32) {

	drawEnemy(screen, sin, cos, 0, ScreenHeight)
	drawEnemy(screen, sin, cos, 0, -ScreenHeight)
	drawEnemy(screen, sin, cos, ScreenWidth, 0)
	drawEnemy(screen, sin, cos, -ScreenWidth, 0)

}

func drawEnemy(screen *ebiten.Image, sin float32, cos float32, addX float32, addY float32) {

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

	for i := 0; i < len(Enemies); i++ {
		enemy := Enemies[i]
		if enemy != nil {
			enemyY := enemy.Y + addY
			enemyX := enemy.X + addX
			drawEnemyImage(screen, EnemyImage, enemyX, enemyY, float32(EnemyBounds.Dx()), float32(EnemyBounds.Dy()), XY.Rotate, playerX, playerY)

		}
	}
}

func ResetEnemies() {
	for i := 0; i < len(Enemies); i++ {
		Enemies[i] = nil
	}
}

func UpdateEnemies() {

	for i := 0; i < len(Enemies); i++ {

		enemy := Enemies[i]

		if enemy == nil {
			continue
		}

		away := CalculateDistanceBetweenPointsWithoutSqrt(enemy.X, enemy.Y, XY.X, XY.Y)

		if away > 60000 {
			Enemies[i] = nil
			continue
		}

		enemy.X = enemy.X + enemy.Dirx
		enemy.Y = enemy.Y + enemy.Diry

		if !enemy.moveAway && Random_range(1, 100) > 75 {
			zeroIn(enemy)
			if away < 4000 {
				enemy.moveAway = true
				enemy.Dirx, enemy.Diry = setDirectionToPlayer(enemy.X, enemy.Y)
				enemy.Dirx = enemy.Dirx * -1
				enemy.Diry = enemy.Diry * -1
			}
		}

		if enemy.X < -HalfSpaceWidth {
			enemy.X = enemy.X + HalfSpaceWidth*2
		}
		if enemy.X > HalfSpaceWidth {
			enemy.X = enemy.X - HalfSpaceWidth*2
		}
		if enemy.Y < -HalfSpaceHeight {
			enemy.Y = enemy.Y + HalfSpaceHeight*2
		}
		if enemy.Y > HalfSpaceHeight {
			enemy.Y = enemy.Y - HalfSpaceHeight*2
		}

		if Random_range(1, 10000) > 9950 {

			theta := 180.0 / math.Pi * math.Atan2(float64(XY.X-enemy.X), float64(XY.Y-enemy.Y))
			FireBullet(float32(ShipBounds.Dy())*1.5, enemy.X, enemy.Y, float32(theta), false)
		}

	}

}

func zeroIn(enemy *Enemy) {

	theta := 180.0 / math.Pi * math.Atan2(float64(XY.X-enemy.X), float64(XY.Y-enemy.Y))

	sin := math.Sin(theta * math.Pi / 180.0)
	cos := math.Cos(theta * math.Pi / 180.0)

	x1 := float64(0.4)
	y1 := float64(0.4)

	enemy.Dirx = float32(x1*cos + y1*sin)
	enemy.Diry = float32(-x1*sin + y1*cos)

}
