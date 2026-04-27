package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	X, Y          float64
	VX, VY        float64
	OnGround      bool
	FacingRight   bool
	AnimFrame     int
	AnimTick      int
	Dead          bool
	Won           bool
}

func NewPlayer() *Player {
	return &Player{
		X:           TileSize * 1.5,
		Y:           TileSize * 14,
		FacingRight: true,
	}
}

func (p *Player) Reset() {
	p.X = TileSize * 1.5
	p.Y = TileSize * 14
	p.VX = 0
	p.VY = 0
	p.OnGround = false
	p.FacingRight = true
	p.AnimFrame = 0
	p.AnimTick = 0
	p.Dead = false
	p.Won = false
}

func (p *Player) Update(level [][]int, coins *[]Coin) bool {
	if p.Dead || p.Won {
		return false
	}

	// Horizontal movement
	p.VX = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.VX = -MoveSpeed
		p.FacingRight = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.VX = MoveSpeed
		p.FacingRight = true
	}

	// Jump
	jumped := ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeySpace)
	if jumped && p.OnGround {
		p.VY = JumpStrength
		p.OnGround = false
	}

	// Gravity
	p.VY += Gravity
	if p.VY > 12 {
		p.VY = 12
	}

	// Move X
	p.X += p.VX
	p.resolveCollisionsX(level)

	// Move Y
	p.Y += p.VY
	p.OnGround = false
	p.resolveCollisionsY(level)

	// Animation
	if p.VX != 0 && p.OnGround {
		p.AnimTick++
		if p.AnimTick >= 8 {
			p.AnimTick = 0
			p.AnimFrame = (p.AnimFrame + 1) % 2
		}
	} else {
		p.AnimFrame = 0
	}

	// Collect coins
	px, py := int(p.X)/TileSize, int(p.Y)/TileSize
	for i := range *coins {
		c := &(*coins)[i]
		if !c.Collected && c.Col >= px-1 && c.Col <= px+1 && c.Row >= py-1 && c.Row <= py+1 {
			cx := float64(c.Col*TileSize) + TileSize/2
			cy := float64(c.Row*TileSize) + TileSize/2
			pcx := p.X + TileSize/2
			pcy := p.Y + TileSize/2
			if abs(cx-pcx) < TileSize && abs(cy-pcy) < TileSize {
				c.Collected = true
			}
		}
	}

	// Check spikes
	col1 := int(p.X+2) / TileSize
	col2 := int(p.X+TileSize-3) / TileSize
	row1 := int(p.Y+2) / TileSize
	row2 := int(p.Y+TileSize-3) / TileSize
	for _, col := range []int{col1, col2} {
		for _, row := range []int{row1, row2} {
			if row >= 0 && row < LevelRows && col >= 0 && col < LevelCols {
				if level[row][col] == TileSpike {
					p.Dead = true
					return false
				}
				if level[row][col] == TileGoal {
					p.Won = true
					return true
				}
			}
		}
	}

	// Fell out of world
	if p.Y > float64(LevelRows*TileSize) {
		p.Dead = true
	}

	return p.Won
}

func (p *Player) resolveCollisionsX(level [][]int) {
	col1 := int(p.X) / TileSize
	col2 := int(p.X+TileSize-1) / TileSize
	row1 := int(p.Y+2) / TileSize
	row2 := int(p.Y+TileSize-1) / TileSize

	for _, row := range []int{row1, row2} {
		for _, col := range []int{col1, col2} {
			if isSolid(level, col, row) {
				if p.VX > 0 {
					p.X = float64(col*TileSize) - TileSize
				} else if p.VX < 0 {
					p.X = float64((col+1) * TileSize)
				}
				p.VX = 0
			}
		}
	}
}

func (p *Player) resolveCollisionsY(level [][]int) {
	col1 := int(p.X+2) / TileSize
	col2 := int(p.X+TileSize-3) / TileSize
	row1 := int(p.Y) / TileSize
	row2 := int(p.Y+TileSize-1) / TileSize

	for _, col := range []int{col1, col2} {
		if isSolid(level, col, row2) && p.VY >= 0 {
			p.Y = float64(row2*TileSize) - TileSize
			p.VY = 0
			p.OnGround = true
		}
		if isSolid(level, col, row1) && p.VY < 0 {
			p.Y = float64((row1+1) * TileSize)
			p.VY = 0
		}
	}
}

func isSolid(level [][]int, col, row int) bool {
	if row < 0 || row >= LevelRows || col < 0 || col >= LevelCols {
		return true
	}
	t := level[row][col]
	return t == TileSolid
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (p *Player) Draw(screen *ebiten.Image) {
	x, y := float64(int(p.X)), float64(int(p.Y))
	sz := float32(TileSize)

	// 1. Ciało Gophera (Niebieski owal/prostokąt)
	// Gopher jest krępy, więc wypełniamy większość Tile'a
	vector.DrawFilledRect(screen, float32(x+2), float32(y+4), sz-4, sz-6, color.RGBA{100, 200, 255, 255}, false)

	// 2. Brzuszek (Jaśniejszy niebieski lub białawy)
	vector.DrawFilledRect(screen, float32(x+4), float32(y+10), sz-8, sz-12, color.RGBA{180, 230, 255, 255}, false)

	// 3. Oczy (Duże białe koła/kwadraty)
	// Przesuwamy oczy w stronę, w którą patrzy Gopher
	eyeOffset := float32(0)
	if p.FacingRight {
		eyeOffset = 2
	} else {
		eyeOffset = -2
	}

	// Lewe oko
	vector.DrawFilledRect(screen, float32(x)+4+eyeOffset, float32(y)+2, 5, 5, color.White, false)
	// Prawe oko
	vector.DrawFilledRect(screen, float32(x)+9+eyeOffset, float32(y)+2, 5, 5, color.White, false)

	// 4. Źrenice (Czarne kropki)
	pupilOffset := float32(0)
	if p.FacingRight {
		pupilOffset = 1.5
	}
	vector.DrawFilledRect(screen, float32(x)+5+eyeOffset+pupilOffset, float32(y)+3.5, 2, 2, color.Black, false)
	vector.DrawFilledRect(screen, float32(x)+10+eyeOffset+pupilOffset, float32(y)+3.5, 2, 2, color.Black, false)

	// 5. Nos (Mały różowy/brązowy kwadracik pod oczami)
	vector.DrawFilledRect(screen, float32(x)+7.5+eyeOffset, float32(y)+6, 3, 2, color.RGBA{255, 150, 150, 255}, false)

	// 6. Zęby (Dwa małe białe prostokąciki)
	vector.DrawFilledRect(screen, float32(x)+7+eyeOffset, float32(y)+8, 1.5, 2, color.White, false)
	vector.DrawFilledRect(screen, float32(x)+9.5+eyeOffset, float32(y)+8, 1.5, 2, color.White, false)

	// 7. Łapki / Nogi (Animizowane)
	legOff := float32(0)
	if p.AnimFrame == 1 {
		legOff = 2
	}
	
	// Kolor łapek (kremowy/beżowy)
	pawColor := color.RGBA{255, 220, 180, 255}
	
	// Lewa noga
    // Najpierw liczymy pozycję na float64, a potem całość zamieniamy na float32
    vector.DrawFilledRect(screen, float32(x+3), float32(y+float64(sz)-4)+legOff, 4, 3, pawColor, false)

    // Prawa noga
    // Tutaj sz musi być rzutowane na float64, żeby dodać się do x
    vector.DrawFilledRect(screen, float32(x+float64(sz)-7), float32(y+float64(sz)-4)-legOff, 4, 3, pawColor, false)
}