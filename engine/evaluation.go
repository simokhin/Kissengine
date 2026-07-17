package engine

var pieceValues = [7]Evaluation{
	Pawn:   100,
	Knight: 300,
	Bishop: 300,
	Rook:   500,
	Queen:  900,
	King:   0,
}

type Evaluation int

func Evaluate(board BoardState) Evaluation {
	var evaluation Evaluation
	phase := gamePhase(board)

	for i := range board.squares {
		piece := board.squares[i]
		if piece == Empty {
			continue
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

	return evaluation
}
