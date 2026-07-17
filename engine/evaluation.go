package engine

var pieceValues = [7]Evaluation{
	Pawn:   100,
	Knight: 300,
	Bishop: 300,
	Rook:   500,
	Queen:  900,
	King:   0,
}

const (
	bishopPairBonus     Evaluation = 30
	doubledPawnPenalty  Evaluation = 15
	isolatedPawnPenalty Evaluation = 15
)

type Evaluation int

func Evaluate(board BoardState) Evaluation {
	var evaluation Evaluation
	phase := gamePhase(board)
	var whiteBishops, blackBishops int
	var whitePawnsByFile [8]int
	var blackPawnsByFile [8]int

	for i := range board.squares {
		file, _ := SquareIndexToFileRank(Square(i))
		piece := board.squares[i]
		if piece == Empty {
			continue
		}

		if piece.Type() == Pawn {
			if piece.Color() == White {
				whitePawnsByFile[file] += 1
			} else {
				blackPawnsByFile[file] += 1
			}
		}

		if piece.Type() == Bishop {
			if piece.Color() == White {
				whiteBishops++
			} else {
				blackBishops++
			}
		}

		var positional Evaluation
		if piece.Type() == King {
			mg := pstValue(KingMiddlegamePST, Square(i), piece.Color())
			eg := pstValue(KingEndGamePST, Square(i), piece.Color())
			positional = Evaluation(phase*float64(mg) + (1-phase)*float64(eg))
		} else {
			positional = pstValue(pstTables[piece.Type()], Square(i), piece.Color())
		}

		if piece.Color() == board.SideToMove().Color() {
			evaluation += pieceValues[piece.Type()] + positional
		} else {
			evaluation -= pieceValues[piece.Type()] + positional
		}
	}

	for i := range board.squares {
		piece := board.squares[i]
		if piece.Type() != Pawn {
			continue
		}

		file, _ := SquareIndexToFileRank(Square(i))

		var pawnsByFile [8]int
		if piece.Color() == White {
			pawnsByFile = whitePawnsByFile
		} else {
			pawnsByFile = blackPawnsByFile
		}

		var penalty Evaluation
		if pawnsByFile[file] > 1 {
			penalty += doubledPawnPenalty
		}

		isolated := true
		if file > 0 && pawnsByFile[file-1] > 0 {
			isolated = false
		}
		if file < 7 && pawnsByFile[file+1] > 0 {
			isolated = false
		}
		if isolated {
			penalty += isolatedPawnPenalty
		}

		if piece.Color() == board.SideToMove().Color() {
			evaluation -= penalty
		} else {
			evaluation += penalty
		}
	}

	if whiteBishops >= 2 {
		if board.SideToMove().Color() == White {
			evaluation += bishopPairBonus
		} else {
			evaluation -= bishopPairBonus
		}
	}
	if blackBishops >= 2 {
		if board.SideToMove().Color() == Black {
			evaluation += bishopPairBonus
		} else {
			evaluation -= bishopPairBonus
		}
	}

	return evaluation
}
