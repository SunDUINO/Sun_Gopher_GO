package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	
)

type Game struct {

	ScreenWidth  int
    ScreenHeight int

	State      GameState
	Player     *Player
	Level      [][]int
	Coins      []Coin
	Enemies    []*Enemy
	Score      int
	Lives      int
	Timer      int
	BgTick     float64
	// Title screen
	TitleAnim  float64
}

func NewGame( sWidth, sHeight int) *Game {
	g := &Game{

		ScreenWidth:  sWidth,  // Zapisujemy przekazaną szerokość
        ScreenHeight: sHeight, // Zapisujemy przekazaną wysokość
		State: StateTitle,
		Lives: 3,
	}
	g.initLevel()
	return g
}

func (g *Game) initLevel() {
	// Deep copy level data so we can modify it
	g.Level = make([][]int, LevelRows)
	for i := range LevelData {
		row := make([]int, LevelCols)
		copy(row, LevelData[i])
		g.Level[i] = row
	}
	g.Player = NewPlayer()
	g.Coins = BuildCoins(g.Level)
	g.Enemies = []*Enemy{
		NewEnemy(12, 15),
		NewEnemy(18, 15),
		NewEnemy(6, 10),
	}
	g.Timer = 0
}

func (g *Game) restartLevel() {
	g.initLevel()
	if g.Lives <= 0 {
		g.State = StateGameOver
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func (g *Game) Update() error {
	g.BgTick += 0.02

	switch g.State {
	case StateTitle:
		g.TitleAnim += 0.05
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State = StatePlaying
		}

	case StatePlaying:
		g.Timer++
		for i := range g.Coins {
			g.Coins[i].Update()
		}
		for _, e := range g.Enemies {
			e.Update(g.Level, g.Player)
		}
		won := g.Player.Update(g.Level, &g.Coins)

		// Count newly collected coins
		for i := range g.Coins {
			if g.Coins[i].Collected {
				g.Score += 10
				g.Coins[i].Collected = true // already set, just ensure
			}
		}

		if won {
			g.State = StateWin
		} else if g.Player.Dead {
			g.Lives--
			if g.Lives <= 0 {
				g.State = StateGameOver
			} else {
				g.Player.Reset()
				g.initLevel()
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.State = StateTitle
		}

	case StateGameOver, StateWin:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Score = 0
			g.Lives = 3
			g.State = StatePlaying
			g.initLevel()
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.State = StateTitle
		}
	}
	return nil
}

// ── Draw ──────────────────────────────────────────────────────────────────────

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case StateTitle:
		g.drawTitle(screen)
	case StatePlaying:
		g.drawGame(screen)
	case StateGameOver:
		g.drawOverlay(screen, "GAME OVER", color.RGBA{180, 20, 20, 255})
	case StateWin:
		g.drawOverlay(screen, "YOU WIN!", color.RGBA{20, 180, 60, 255})
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	// Gradient background
	for y := 0; y < int(g.ScreenHeight); y++ {
		t := float64(y) / float64(int(g.ScreenHeight))
		r := uint8(20 + t*60)
		b := uint8(80 + t*80)
		vector.DrawFilledRect(screen, 0, float32(y), float32(g.ScreenWidth), 1,
			color.RGBA{r, 10, b, 255}, false)
	}

	// Animated stars
	for i := 0; i < 40; i++ {
		sx := float32(math.Mod(float64(i*37+47), float64(g.ScreenWidth)))
		sy := float32(math.Mod(float64(i*53+17), float64(g.ScreenHeight)))
		br := uint8(150 + uint8(math.Sin(g.TitleAnim+float64(i))*50))
		vector.DrawFilledCircle(screen, sx, sy, 1, color.RGBA{br, br, br, 255}, false)
	}

	// Platforms deco
	for i := 0; i < 5; i++ {
		px := float32(i*100) + float32(math.Sin(g.TitleAnim+float64(i))*8)
		py := float32(180 + i*10)
		vector.DrawFilledRect(screen, px, py, 80, 8, color.RGBA{100, 60, 20, 255}, false)
		vector.DrawFilledRect(screen, px, py-2, 80, 4, color.RGBA{60, 140, 60, 255}, false)
	}

	// Title text
	ebitenutil.DebugPrintAt(screen, "G O P H E R   G O !", 140, 60)
	ebitenutil.DebugPrintAt(screen, "~ SunRiver platformer ~", 148, 80)
	ebitenutil.DebugPrintAt(screen, "Arrow / WASD  – move", 150, 130)
	ebitenutil.DebugPrintAt(screen, "Up / W / Space – jump", 148, 145)
	ebitenutil.DebugPrintAt(screen, "Stomp enemies!", 168, 160)
	if int(g.TitleAnim*3)%2 == 0 {
		ebitenutil.DebugPrintAt(screen, ">> PRESS ENTER TO START <<", 120, 210)
	}
}

func (g *Game) drawGame(screen *ebiten.Image) {
	// Sky gradient
	for y := 0; y < int(g.ScreenHeight); y++ {
		t := float64(y) / float64(int(g.ScreenHeight))
		r := uint8(80 + t*60)
		gg := uint8(140 + t*60)
		b := uint8(220 - t*60)
		vector.DrawFilledRect(screen, 0, float32(y), float32(g.ScreenWidth), 1,
			color.RGBA{r, gg, b, 255}, false)
	}

	// Draw tiles
	for row := 0; row < LevelRows; row++ {
		for col := 0; col < LevelCols; col++ {
			tile := g.Level[row][col]
			x := float32(col * TileSize)
			y := float32(row * TileSize)
			sz := float32(TileSize)

			switch tile {
			case TileSolid:
				// Grass top or dirt
				vector.DrawFilledRect(screen, x, y, sz, sz, color.RGBA{120, 80, 40, 255}, false)
				// Grass on top
				if row > 0 && g.Level[row-1][col] != TileSolid {
					vector.DrawFilledRect(screen, x, y, sz, 4, color.RGBA{60, 160, 60, 255}, false)
				}
				// Brick lines
				vector.StrokeRect(screen, x+0.5, y+0.5, sz-1, sz-1, 1,
					color.RGBA{90, 55, 25, 255}, false)
			case TileSpike:
				// Triangle spikes
				for s := 0; s < 2; s++ {
					sx := x + float32(s)*8 + 2
					vector.DrawFilledRect(screen, sx, y+sz-5, 6, 5,
						color.RGBA{180, 30, 30, 255}, false)
					// tip
					vector.DrawFilledRect(screen, sx+2, y+sz-10, 2, 8,
						color.RGBA{220, 50, 50, 255}, false)
				}
			case TileGoal:
				// Goal star / portal
				blink := uint8(150 + uint8(math.Sin(float64(g.Timer)*0.1)*80))
				vector.DrawFilledRect(screen, x+2, y+2, sz-4, sz-4,
					color.RGBA{blink, blink, 50, 255}, false)
				vector.DrawFilledCircle(screen, x+sz/2, y+sz/2, 5,
					color.RGBA{255, 255, 50, 255}, false)
				ebitenutil.DebugPrintAt(screen, "G", int(x)+5, int(y)+4)
			}
		}
	}

	// Draw coins
	for i := range g.Coins {
		g.Coins[i].Draw(screen)
	}

	// Draw enemies
	for _, e := range g.Enemies {
		e.Draw(screen)
	}

	// Draw player
	g.Player.Draw(screen)

	// HUD
	collected := 0
	for _, c := range g.Coins {
		if c.Collected {
			collected++
		}
	}
	hud := fmt.Sprintf("LIVES:%d  SCORE:%d  COINS:%d/%d  TIME:%d",
		g.Lives, g.Score, collected, len(g.Coins), g.Timer/60)
	// HUD bar
	vector.DrawFilledRect(screen, 0, 0, float32(g.ScreenWidth), 10, color.RGBA{0, 0, 0, 160}, false)
	ebitenutil.DebugPrintAt(screen, hud, 2, 1)
}

func (g *Game) drawOverlay(screen *ebiten.Image, msg string, _ color.RGBA) {
	g.drawGame(screen)
	vector.DrawFilledRect(screen, 0, 0, float32(g.ScreenWidth), float32(g.ScreenHeight),
		color.RGBA{0, 0, 0, 160}, false)
	x := int(g.ScreenWidth)/2 - len(msg)*3
	ebitenutil.DebugPrintAt(screen, msg, x, int(g.ScreenHeight)/2-10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.Score), x+4, int(g.ScreenHeight)/2+5)
	ebitenutil.DebugPrintAt(screen, "ENTER – restart   ESC – menu", x-30, int(g.ScreenHeight)/2+20)
	
}

// ── Layout ────────────────────────────────────────────────────────────────────

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return int(g.ScreenWidth), int(g.ScreenHeight)
}
