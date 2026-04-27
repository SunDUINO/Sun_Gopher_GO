/*
 * ╔════════════════════════════════════════════════════════════════╗
 * ║ sun_gopher_go                                                  ║
 * ║ Plik / File: main.go                                           ║
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

