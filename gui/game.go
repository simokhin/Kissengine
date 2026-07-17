package main

import (
	"MyChessEngine/engine"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	board             engine.BoardState
	humanColor        engine.Piece
	hasSelected       bool
	selected          engine.Square
	legalMoves        []engine.Move
	engineThinking    bool
	engineResult      chan engine.SearchResult
	awaitingPromotion bool
	promotionMoves    []engine.Move
	gameOver          bool
	gameOverText      string
}

func (g *Game) applyMove(move engine.Move) {
	g.board = engine.MakeMove(g.board, move)

	sound := moveSound(move, g.board)
	sounds[sound].Rewind()
	sounds[sound].Play()

	if len(engine.GenerateLegalMoves(g.board)) == 0 {
		g.gameOver = true
		if g.board.InCheck() {
			if g.board.SideToMove().Color() == engine.White {
				g.gameOverText = "Checkmate! Black wins"
			} else {
				g.gameOverText = "Checkmate! White wins"
			}
		} else {
			g.gameOverText = "Stalemate - draw"
		}
		return
	}

	if g.board.SideToMove().Color() != g.humanColor {
		g.engineThinking = true
		g.engineResult = make(chan engine.SearchResult, 1)
		board := g.board
		go func() {
			g.engineResult <- engine.FindBestMoveByTime(board, time.Duration(moveTime)*time.Millisecond, nil, false)
		}()
	}
}

const (
	squareSize   = 100
	screenWidth  = squareSize * 8
	screenHeight = squareSize * 8

	searchDepth = 5
	moveTime    = 1000
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

func moveSound(move engine.Move, boardAfter engine.BoardState) SoundType {
	if boardAfter.InCheck() {
		return CheckSound
	}
	if move.Flag() == engine.Capture || move.Flag() == engine.EnPassantCapture {
		return CaptureSound
	}
	if move.Flag() == engine.KingCastle || move.Flag() == engine.QueenCastle {
		return CastleSound
	}
	if move.Promotion() != engine.Empty {
		return PromoteSound
	}
	return MoveSound
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	if g.engineThinking {
		select {
		case result := <-g.engineResult:
			g.engineThinking = false
			g.applyMove(result.Move)
		default:
			return nil
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		file := x / squareSize
		rank := 7 - y/squareSize
		clicked := engine.FileRankToSquareIndex(file, rank)

		if g.awaitingPromotion {
			destSquare := g.promotionMoves[0].To()
			destFile, destRank := engine.SquareIndexToFileRank(destSquare)
			step := -1
			if destRank == 0 {
				step = 1
			}
			promotionPieces := []engine.Piece{engine.Queen, engine.Rook, engine.Bishop, engine.Knight}
			for i, p := range promotionPieces {
				slotSquare := engine.FileRankToSquareIndex(destFile, destRank+i*step)
				if clicked == slotSquare {
					for _, m := range g.promotionMoves {
						if m.Promotion() == p {
							g.applyMove(m)
							break
						}
					}
					break
				}
			}

			g.awaitingPromotion = false
			g.promotionMoves = nil
			g.hasSelected = false
			g.legalMoves = nil
			return nil
		}

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
			var candidateMoves []engine.Move
			for _, m := range g.legalMoves {
				if m.To() == clicked {
					candidateMoves = append(candidateMoves, m)
				}
			}

			switch len(candidateMoves) {
			case 0:
				// клик мимо легального хода — просто снимаем выбор ниже
			case 1:
				g.applyMove(candidateMoves[0])
			default:
				g.awaitingPromotion = true
				g.promotionMoves = candidateMoves
			}

			if !g.awaitingPromotion {
				g.hasSelected = false
				g.legalMoves = nil
			}
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

			// Draw pieces
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

	// Draw transparent circles in the squares of possible moves of the piece
	if g.hasSelected {
		for _, move := range g.legalMoves {
			toFile, toRank := engine.SquareIndexToFileRank(move.To())
			cx := float32(toFile*squareSize) + squareSize/2
			cy := float32((7-toRank)*squareSize) + squareSize/2

			isCapture := move.Flag() == engine.Capture || move.Flag() == engine.EnPassantCapture
			if isCapture {
				vector.StrokeCircle(screen, cx, cy, squareSize/2-4, 4, color.RGBA{0, 0, 0, 120}, true)
			} else {
				vector.FillCircle(screen, cx, cy, squareSize/6, color.RGBA{0, 0, 0, 80}, true)
			}
		}
	}

	if g.awaitingPromotion {
		destSquare := g.promotionMoves[0].To()
		destFile, destRank := engine.SquareIndexToFileRank(destSquare)

		step := -1
		if destRank == 0 {
			step = 1
		}

		promotionPieces := []engine.Piece{engine.Queen, engine.Rook, engine.Bishop, engine.Knight}

		for i, p := range promotionPieces {
			rank := destRank + i*step
			x := float32(destFile * squareSize)
			y := float32((7 - rank) * squareSize)

			vector.FillRect(screen, x, y, squareSize, squareSize, color.RGBA{200, 200, 200, 255}, false)

			img := pieceImages[g.humanColor|p]
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

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, g.gameOverText, screenWidth/2-len(g.gameOverText)*3, screenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
