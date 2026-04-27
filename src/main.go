package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"sun_gopher_go/src/game"
)

const (
	ScreenWidth  = 480
	ScreenHeight = 270
	GameTitle    = "Gopher Go!"
)

func main() {
	ebiten.SetWindowSize(ScreenWidth*3, ScreenHeight*3)
	ebiten.SetWindowTitle(GameTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame(ScreenWidth, ScreenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

