package main

import (
	"MyChessEngine/engine"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	board       engine.BoardState
	humanColor  engine.Piece
	hasSelected bool
	selected    engine.Square
	legalMoves  []engine.Move
}

const (
	squareSize   = 100
	screenWidth  = squareSize * 8
	screenHeight = squareSize * 8

	searchDepth = 5
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		file := x / squareSize
		rank := 7 - y/squareSize
		clicked := engine.FileRankToSquareIndex(file, rank)

		if !g.hasSelected {
			var movesFromClicked []engine.Move
			for _, m := range engine.GenerateLegalMoves(g.board) {
				if m.From() == clicked {
					movesFromClicked = append(movesFromClicked, m)
				}
			}
			if len(movesFromClicked) > 0 {
				g.selected = clicked
				g.hasSelected = true
				g.legalMoves = movesFromClicked
			}
			fmt.Println(g.selected)
			fmt.Println(len(movesFromClicked))
		} else {
			var move engine.Move
			found := false
			for _, m := range g.legalMoves {
				if m.To() == clicked {
					move = m
					found = true
					break
				}
			}

			if found {
				g.board = engine.MakeMove(g.board, move)

				if g.board.SideToMove().Color() != g.humanColor {
					engineMove := engine.FindBestMove(g.board, searchDepth)
					g.board = engine.MakeMove(g.board, engineMove)
				}
			}

			g.hasSelected = false
			g.legalMoves = nil
		}
	}
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
				img := pieceImages[piece]
				if img != nil {
					scale := float64(squareSize) / 256.0
					op := &ebiten.DrawImageOptions{}
					op.Filter = ebiten.FilterLinear
					op.GeoM.Scale(scale, scale)
					op.GeoM.Translate(float64(x), float64(y))
					screen.DrawImage(img, op)
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
