package main

import (
	"MyChessEngine/engine"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Kissengine GUI")

	// testFEN := "k7/7P/8/3q4/8/5n2/8/7K w - - 0 1"

	game := &Game{board: engine.ParseFEN(engine.StartFen), humanColor: engine.White}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
