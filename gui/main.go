package main

import (
	"MyChessEngine/engine"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Kissengine GUI")

	testFen := "k7/5P2/8/8/8/8/8/K7 w - - 0 1"

	game := &Game{board: engine.ParseFEN(testFen), humanColor: engine.White}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
