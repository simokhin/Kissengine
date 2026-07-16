package main

import (
	"MyChessEngine/engine"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Kissengine GUI")

	game := &Game{board: engine.ParseFEN(engine.StartFen)}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
