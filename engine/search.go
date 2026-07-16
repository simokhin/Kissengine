package engine

const (
	Infinity Evaluation = 1_000_000
	Mate     Evaluation = 100_000
)

func FindBestMove(board BoardState, depth int) Move {
	var bestMove Move

	moves := GenerateLegalMoves(board)

	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		moveEvaluation := -NegaMax(newBoard, depth-1)
		if moveEvaluation > best {
			best = moveEvaluation
			bestMove = move
		}
	}

	return bestMove
}

func NegaMax(board BoardState, depth int) Evaluation {
	if depth == 0 {
		return Evaluate(board)
	}

	moves := GenerateLegalMoves(board)
	if len(moves) == 0 {
		if board.InCheck() {
			return -(Mate - Evaluation(depth))
		}
		return 0
	}
	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score := -NegaMax(newBoard, depth-1)
		if score > best {
			best = score
		}
	}

	return best
}
