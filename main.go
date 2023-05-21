package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"bjason.org/goroid/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	noiseImage    *image.RGBA
	musicPlayer   *game.Player
	keyboard      game.Keyboard
	Level         int
	finished      bool
	messageScroll int
}

const ROCK_INITAL_MINIMUM = 70

func initialiseLandscape() {

	game.LoadImages()

	game.CentreRockPolyArrayAndFlipCoOrds()

	initialiseRocks()

}

func initialiseRocks() {
	rand.Seed(time.Now().UnixNano())
	rockCount := 0
	i := 0
	for y := -game.HalfSpaceHeight + 64; y < game.HalfSpaceHeight-64; y += 64 {

		for x := -game.HalfSpaceWidth + 64; x < game.HalfSpaceWidth-64; x += 64 {

			away := game.CalculateDistanceBetweenPointsWithoutSqrt(float32(x), float32(y), game.XY.X, game.XY.Y)

			if away > 36000 {

				random := rand.Intn(80)

				if random > 60 {
					for i < len(game.Rocks) {
						if game.Rocks[i] == nil {
							break
						}
						i++
					}
					if i >= len(game.Rocks) {
						continue
					}

					xx := x
					yy := y
					size := ROCK_INITAL_MINIMUM

					game.Rocks[i] = game.NewRock(float32(xx), float32(yy), float32(size))
					i++
					rockCount++

				}
			}
		}
	}
	log.Printf("************ ADDED ROCKS %d\n", rockCount)
}

func countRockstoChangeLevel() bool {
	count := 0

	for i := 0; i < len(game.Rocks); i++ {
		if game.Rocks[i] != nil && game.Rocks[i].Scale >= ROCK_INITAL_MINIMUM {
			count++
		}
	}
	return count <= 2
}

func (g *Game) Update() error {

	g.handleKeyboard()

	if !g.finished {
		g.playingUpdateLogic()
	}

	game.UpdateParticles()

	if g.keyboard.Exit && runtime.GOARCH != "js" && runtime.GOOS != "js" {
		return ebiten.Termination
	}

	if game.XY.Lives <= 0 {
		if !g.finished {
			g.keyboard.Space = false
		}
		g.finished = true

	}
	if g.finished && g.keyboard.Space {
		g.resetTheGame()
	}
	return nil
}

func (g *Game) resetTheGame() {
	g.keyboard.Space = false
	game.XY.Lives = game.LIVES
	g.Level = 1
	game.XY.Score = 0
	game.ResetBullets()
	game.ResetRocks()
	game.ResetEnemies()
	g.finished = false
}

func (g *Game) playingUpdateLogic() {
	if countRockstoChangeLevel() {
		initialiseRocks()
		g.Level++
	}

	game.UpdateRock()

	game.PlayerRotateAndMoveShip(g.keyboard, g.musicPlayer)
	game.CheckRockHit(g.musicPlayer)
	game.EnemyHitbyBullet()

	game.UpdateBullets()

	if game.Random_range(1, 100) > 95 {
		topBottom := float32(game.Random_range(-1, 2))
		leftRight := float32(game.Random_range(-1, 2))
		if topBottom != 0 && leftRight != 0 {
			game.AddEnemy(game.XY.X-(game.ScreenWidth/2*leftRight), game.XY.Y+(game.ScreenHeight/3*topBottom))
		}

	}

	game.UpdateEnemies()

	g.musicPlayer.SoundReadyToPlay()
}

func (g *Game) handleKeyboard() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.keyboard.Left = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		g.keyboard.Left = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.keyboard.Right = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		g.keyboard.Right = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.keyboard.Space = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.keyboard.Space = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.keyboard.Exit = true
	}

	/*
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.keyboard.Down = true
		} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
			g.keyboard.Down = false
		}
	*/
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.keyboard.Up = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		g.keyboard.Up = false
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	sin := float32(math.Sin(float64(-game.XY.Rotate) * math.Pi / 180.0))
	cos := float32(math.Cos(float64(-game.XY.Rotate) * math.Pi / 180.0))

	game.DrawStars(screen, sin, cos)
	game.DrawRocks(screen)

	game.DrawEnemy(screen, sin, cos)

	game.DrawBullets(screen, sin, cos)

	game.DrawShip(screen, game.PLAYER_SCREEN_X, game.PLAYER_SCREEN_Y+1, 0)

	game.DrawExplosions(screen, sin, cos)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("score=%d lives=%d level=%d",
		game.XY.Score, game.XY.Lives, g.Level))

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("tps: %0.2f fps: %0.2f",
		ebiten.ActualTPS(), ebiten.ActualFPS()), 0, game.ScreenHeight-16)

	if g.finished {
		ebitenutil.DebugPrintAt(screen, "Game over. Press space to play again. Game over", g.messageScroll/2-300, game.ScreenHeight/2)
		g.messageScroll++
		g.messageScroll = g.messageScroll % (game.ScreenWidth * 3)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.ScreenWidth, game.ScreenHeight
}

func main() {
	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		ebiten.SetFullscreen(true)
	} else {
		ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
		ebiten.SetWindowTitle("goroid")
	}

	game.SetupStars()
	initialiseLandscape()

	g := &Game{
		noiseImage:  image.NewRGBA(image.Rect(0, 0, game.ScreenWidth, game.ScreenHeight)),
		musicPlayer: game.InitialiseSound(),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
