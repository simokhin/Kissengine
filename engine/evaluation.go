package engine

var pieceValues = [7]Evaluation{
	Pawn:   100,
	Knight: 300,
	Bishop: 300,
	Rook:   500,
	Queen:  900,
	King:   0,
}

const bishopPairBonus Evaluation = 30

type Evaluation int

func Evaluate(board BoardState) Evaluation {
	var evaluation Evaluation
	phase := gamePhase(board)
	var whiteBishops, blackBishops int

	for i := range board.squares {
		piece := board.squares[i]
		if piece == Empty {
			continue
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
