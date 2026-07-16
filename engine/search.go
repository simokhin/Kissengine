package engine

import (
	"sort"
	"time"
)

const (
	Infinity Evaluation = 1_000_000
	Mate     Evaluation = 100_000
)

func moveScore(board BoardState, move Move) int {
	if move.CapturedPiece() == Empty {
		return 0
	}
	attacker := board.PieceAt(move.From())
	return int(pieceValues[move.CapturedPiece().Type()])*10 - int(pieceValues[attacker.Type()])
}

func orderMoves(board BoardState, moves []Move) []Move {
	sort.Slice(moves, func(i, j int) bool {
		return moveScore(board, moves[i]) > moveScore(board, moves[j])
	})
	return moves
}

func FindBestMoveByTime(board BoardState, timeLimit time.Duration) Move {
	deadline := time.Now().Add(timeLimit)
	var bestMove Move

	for depth := 1; ; depth++ {
		if time.Now().After(deadline) {
			break
		}
		bestMove = FindBestMove(board, depth)
	}

	return bestMove
}

func FindBestMove(board BoardState, depth int) Move {
	var bestMove Move

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves)

	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		moveEvaluation := -NegaMax(newBoard, depth-1, -Infinity, -best)
		if moveEvaluation > best {
			best = moveEvaluation
			bestMove = move
		}
	}

	return bestMove
}

func NegaMax(board BoardState, depth int, alpha, beta Evaluation) Evaluation {
	if depth == 0 && !board.InCheck() {
		return Evaluate(board)
	}

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves)

	if len(moves) == 0 {
		if board.InCheck() {
			return -(Mate + Evaluation(depth))
		}
		return 0
	}

	if depth == 0 {
		return Evaluate(board)
	}

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score := -NegaMax(newBoard, depth-1, -beta, -alpha)

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
