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

	for i := range board.squares {
		piece := board.squares[i]
		if piece == Empty || piece.Type() == King {
			continue
		}

		if piece.Color() == board.SideToMove().Color() {
			evaluation += pieceValues[piece.Type()]
		} else {
			evaluation -= pieceValues[piece.Type()]
		}
	}

	return evaluation
}
