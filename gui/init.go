package main

import (
	"MyChessEngine/engine"
	"bytes"
	"embed"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/pieces/*.png
var pieceFS embed.FS

var pieceImages [15]*ebiten.Image

func init() {
	load := func(piece engine.Piece, filename string) {
		data, err := pieceFS.ReadFile("assets/pieces/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}

		pieceImages[piece] = ebiten.NewImageFromImage(img)
	}

	load(engine.White|engine.Pawn, "Chess_plt45.png")
	load(engine.White|engine.Knight, "Chess_nlt45.png")
	load(engine.White|engine.Bishop, "Chess_blt45.png")
	load(engine.White|engine.Rook, "Chess_rlt45.png")
	load(engine.White|engine.Queen, "Chess_qlt45.png")
	load(engine.White|engine.King, "Chess_klt45.png")
	load(engine.Black|engine.Pawn, "Chess_pdt45.png")
	load(engine.Black|engine.Knight, "Chess_ndt45.png")
	load(engine.Black|engine.Bishop, "Chess_bdt45.png")
	load(engine.Black|engine.Rook, "Chess_rdt45.png")
	load(engine.Black|engine.Queen, "Chess_qdt45.png")
	load(engine.Black|engine.King, "Chess_kdt45.png")
}
