package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type PlayerInformation struct {
	X, Y       float32
	Rotate     int
	Score      int
	Lives      int
	NextBullet float32
}

type Keyboard struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
	Exit  bool
	Space bool
}

const LIVES = 10

var XY = PlayerInformation{-0, 0, 0, 0, LIVES, 0}

var Ship_poly_original = []float64{
	0, 6,
	6, 0,
	9, 0,
	15, 6,
	15, 15,
	0, 15,
	0, 6,
}

func CentreShipAbsolute() []float64 {
	var poly []float64 = make([]float64, len(Ship_poly_original))
	copy(poly, Ship_poly_original)
	sin := float32(math.Sin(float64(XY.Rotate) * math.Pi / 180.0))
	cos := float32(math.Cos(float64(XY.Rotate) * math.Pi / 180.0))
	for i := 0; i < len(poly); i += 2 {
		x1 := float32(poly[i] - float64(ShipBounds.Dx())/2)
		y1 := float32(poly[i+1] - float64(ShipBounds.Dy())/2)

		poly[i] = float64(x1*cos + y1*sin)
		poly[i+1] = float64(-x1*sin + y1*cos)

		poly[i] = float64(XY.X) + poly[i]
		poly[i+1] = float64(XY.Y) - poly[i+1]
	}
	return poly
}

func PlayerRotateAndMoveShip(keyboard Keyboard, sound *Player) {

	XY.NextBullet = XY.NextBullet - 1
	if keyboard.Up {

		addX := math.Sin(float64(XY.Rotate) * math.Pi / 180.0)
		addY := math.Cos(float64(XY.Rotate) * math.Pi / 180.0)
		XY.Y = XY.Y + float32(addY)*1.0
		XY.X = XY.X + float32(addX)*1.0
	} else if keyboard.Down {
		addX := math.Sin(float64(XY.Rotate) * math.Pi / 180.0)
		addY := math.Cos(float64(XY.Rotate) * math.Pi / 180.0)
		XY.Y = XY.Y + float32(addY)*-1.0
		XY.X = XY.X + float32(addX)*-1.0
	}
	if keyboard.Right {
		XY.Rotate = XY.Rotate + 2
	} else if keyboard.Left {
		XY.Rotate = XY.Rotate - 2
	}
	if XY.Rotate < 0 {
		XY.Rotate = XY.Rotate + 360
	}
	if XY.Rotate >= 360 {
		XY.Rotate = XY.Rotate - 360
	}

	if keyboard.Space && XY.NextBullet < 0 {
		FireBullet(float32(ShipBounds.Dy())*0.9, XY.X, XY.Y, float32(XY.Rotate), true)
		XY.NextBullet = 10
		sound.playShot()

	}

	wrapScreen()

}

func wrapScreen() {
	if XY.X < -HalfSpaceWidth {
		XY.X = XY.X + HalfSpaceWidth*2
	} else if XY.X > HalfSpaceWidth {
		XY.X = XY.X - HalfSpaceWidth*2
	}

	if XY.Y < -HalfSpaceHeight {
		XY.Y = XY.Y + HalfSpaceHeight*2
	} else if XY.Y > HalfSpaceHeight {
		XY.Y = XY.Y - HalfSpaceHeight*2
	}
}

func DrawShip(screen *ebiten.Image, x float32, y float32, rotate int) {

	s := ShipImage.Bounds().Size()
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)

	op.GeoM.Rotate(float64(rotate%360) * 2 * math.Pi / 360)

	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(ShipImage, op)
}
