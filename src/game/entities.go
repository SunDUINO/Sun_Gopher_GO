/*
 * ╔════════════════════════════════════════════════════════════════╗
 * ║ sun_gopher_go                                                  ║
 * ║ Plik / File: entities.go                                       ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ Autor / Author:                                                ║
 * ║   SunRiver                                                     ║
 * ║   Lothar TeaM                                                  ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ GitHub  : sun_gopher_go                                        ║
 * ║ WWW     : https://lothar-team.pl                               ║
 * ║ Forum   : https://forum.lothar-team.pl                         ║
 * ║                                                                ║
 * ║ Licencja / License: MIT                                        ║
 * ║ Rok / Year: 2026                                               ║
 * ╚════════════════════════════════════════════════════════════════╝
 */
package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ── Coin ──────────────────────────────────────────────────────────────────────

type Coin struct {
	Col, Row  int
	Collected bool
	Anim      float64
}

func BuildCoins(level [][]int) []Coin {
	var coins []Coin
	for row := 0; row < LevelRows; row++ {
		for col := 0; col < LevelCols; col++ {
			if level[row][col] == TileCoin {
				coins = append(coins, Coin{Col: col, Row: row})
			}
		}
	}
	return coins
}

func (c *Coin) Update() {
	c.Anim += 0.12
}

func (c *Coin) Draw(screen *ebiten.Image) {
	if c.Collected {
		return
	}
	cx := float32(c.Col*TileSize) + TileSize/2
	cy := float32(c.Row*TileSize) + TileSize/2 + float32(math.Sin(c.Anim)*2)
	vector.DrawFilledCircle(screen, cx, cy, 5, color.RGBA{255, 220, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx, cy, 3, color.RGBA{255, 255, 150, 255}, false)
}

// ── Enemy (Gumba-style) ───────────────────────────────────────────────────────

type Enemy struct {
	X, Y      float64
	VX        float64
	OnGround  bool
	AnimTick  int
	AnimFrame int
	Dead      bool
}

func NewEnemy(col, row int) *Enemy {
	return &Enemy{
		X:  float64(col * TileSize),
		Y:  float64(row * TileSize),
		VX: -1.0,
	}
}

func (e *Enemy) Update(level [][]int, p *Player) {
	if e.Dead {
		return
	}

	// Gravity
	e.Y += 0.0
	//e.VX = e.VX // keep
	gravity := 0.4
	vy := gravity

	// simple fall
	e.Y += vy
	row := int(e.Y+TileSize) / TileSize
	col1 := int(e.X+2) / TileSize
	col2 := int(e.X+TileSize-3) / TileSize
	for _, col := range []int{col1, col2} {
		if isSolid(level, col, row) {
			e.Y = float64(row*TileSize) - TileSize
			e.OnGround = true
		}
	}

	// Move horizontally
	e.X += e.VX
	colAhead := int(e.X+TileSize+e.VX) / TileSize
	if e.VX < 0 {
		colAhead = int(e.X+e.VX) / TileSize
	}
	midRow := int(e.Y+TileSize/2) / TileSize
	if isSolid(level, colAhead, midRow) || colAhead < 0 || colAhead >= LevelCols {
		e.VX = -e.VX
	}

	// Animation
	e.AnimTick++
	if e.AnimTick >= 10 {
		e.AnimTick = 0
		e.AnimFrame = (e.AnimFrame + 1) % 2
	}

	// Kill player (stomp from above kills enemy, side kills player)
	px, py := p.X, p.Y
	ex, ey := e.X, e.Y
	if !p.Dead && !p.Won &&
		px+TileSize-2 > ex+2 && px+2 < ex+TileSize-2 &&
		py+TileSize-2 > ey+2 && py+2 < ey+TileSize-2 {
		// Check if player is falling onto enemy
		if p.VY > 0 && py+TileSize < ey+TileSize/2+4 {
			e.Dead = true
			p.VY = JumpStrength * 0.6 // bounce
		} else {
			p.Dead = true
		}
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	if e.Dead {
		return
	}
	x := float32(e.X)
	y := float32(e.Y)
	sz := float32(TileSize)

	// Body – brown mushroom shape
	vector.DrawFilledRect(screen, x+2, y+sz/2, sz-4, sz/2, color.RGBA{160, 80, 20, 255}, false)
	// Head cap
	vector.DrawFilledCircle(screen, x+sz/2, y+sz/2, sz/2-1, color.RGBA{180, 50, 10, 255}, false)
	// Eyes
	eyeY := y + sz/2 - 1
	legOff := float32(0)
	if e.AnimFrame == 1 {
		legOff = 1
	}
	vector.DrawFilledRect(screen, x+3, eyeY, 3, 3, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+sz-6, eyeY, 3, 3, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+4, eyeY+1, 2, 2, color.RGBA{10, 10, 10, 255}, false)
	vector.DrawFilledRect(screen, x+sz-5, eyeY+1, 2, 2, color.RGBA{10, 10, 10, 255}, false)
	// Feet
	vector.DrawFilledRect(screen, x+1, y+sz-4+legOff, 5, 4, color.RGBA{100, 50, 10, 255}, false)
	vector.DrawFilledRect(screen, x+sz-6, y+sz-4-legOff, 5, 4, color.RGBA{100, 50, 10, 255}, false)
}
