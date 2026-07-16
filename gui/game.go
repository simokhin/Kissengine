package main

import (
	"MyChessEngine/engine"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	board       engine.BoardState
	hasSelected bool
	selected    engine.Square
	legalMoves  []engine.Move
}

const (
	squareSize   = 64
	screenWidth  = squareSize * 8
	screenHeight = squareSize * 8
)

var (
	lightColor = color.RGBA{238, 238, 210, 1.0}
	darkColor  = color.RGBA{118, 150, 86, 1.0}
)

var pieceGlyphs = [15]string{
	engine.White | engine.Pawn:   "P",
	engine.White | engine.Knight: "N",
	engine.White | engine.Bishop: "B",
	engine.White | engine.Rook:   "R",
	engine.White | engine.Queen:  "Q",
	engine.White | engine.King:   "K",
	engine.Black | engine.Pawn:   "p",
	engine.Black | engine.Knight: "n",
	engine.Black | engine.Bishop: "b",
	engine.Black | engine.Rook:   "r",
	engine.Black | engine.Queen:  "q",
	engine.Black | engine.King:   "k",
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for rank := 0; rank <= engine.Rank8; rank++ {
		for file := 0; file <= engine.FileH; file++ {
			var squareColor color.Color
			if (file+rank)%2 == 0 {
				squareColor = darkColor
			} else {
				squareColor = lightColor
			}

			x := float32(file * squareSize)
			y := float32((7 - rank) * squareSize)

			vector.FillRect(screen, x, y, squareSize, squareSize, squareColor, false)

			square := engine.FileRankToSquareIndex(file, rank)
			piece := g.board.PieceAt(square)

			glyph := pieceGlyphs[piece]
			if glyph != "" {
				ebitenutil.DebugPrintAt(screen, glyph, int(x), int(y))
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
