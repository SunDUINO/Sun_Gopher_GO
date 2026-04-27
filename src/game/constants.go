package game

// Tile types
const (
	TileEmpty    = 0
	TileSolid    = 1
	TileCoin     = 2
	TileSpike    = 3
	TileGoal     = 4

	TileSize     = 16
	Gravity      = 0.4
	JumpStrength = -7.5
	MoveSpeed    = 2.0
)

// State machine
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateGameOver
	StateWin
)
