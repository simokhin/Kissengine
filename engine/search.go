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

func FindBestMove(board BoardState, depth int) Move {
	move, _ := findBestMove(board, depth, time.Time{})
	return move
}

func FindBestMoveByTime(board BoardState, timeLimit time.Duration) Move {
	deadline := time.Now().Add(timeLimit)
	var bestMove Move

	for depth := 1; ; depth++ {
		if time.Now().After(deadline) {
			break
		}
		move, ok := findBestMove(board, depth, deadline)
		if !ok {
			break
		}
		bestMove = move
	}

	return bestMove
}

func findBestMove(board BoardState, depth int, deadline time.Time) (Move, bool) {
	var bestMove Move

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves)

	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, -Infinity, -best, deadline)
		if !ok {
			return bestMove, false
		}
		moveEvaluation := -score
		if moveEvaluation > best {
			best = moveEvaluation
			bestMove = move
		}
	}

	return bestMove, true
}

func quiescence(board BoardState, alpha, beta Evaluation, deadline time.Time) (Evaluation, bool) {
	if !deadline.IsZero() && time.Now().After(deadline) {
		return 0, false
	}

	inCheck := board.InCheck()
	var moves []Move

	if inCheck {
		moves = GenerateLegalMoves(board)
	} else {
		standPat := Evaluate(board)
		if standPat >= beta {
			return beta, true
		}
		if standPat > alpha {
			alpha = standPat
		}
		for _, m := range GenerateLegalMoves(board) {
			if m.CapturedPiece() != Empty {
				moves = append(moves, m)
			}
		}
	}
	orderMoves(board, moves)

	if inCheck && len(moves) == 0 {
		return -Mate, true
	}

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := quiescence(newBoard, -beta, -alpha, deadline)
		if !ok {
			return 0, false
		}
		score = -score

		if score >= beta {
			return beta, true
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha, true
}

func negaMax(board BoardState, depth int, alpha, beta Evaluation, deadline time.Time) (Evaluation, bool) {
	if !deadline.IsZero() && time.Now().After(deadline) {
		return 0, false
	}

	if depth == 0 && !board.InCheck() {
		return quiescence(board, alpha, beta, deadline)
	}

	var bestMove Move

	hash := ComputeHash(board)
	entry, found := Probe(hash)
	if found && entry.depth >= depth {
		switch entry.flag {
		case Exact:
			return entry.evaluation, true
		case LowerBound:
			if entry.evaluation > alpha {
				alpha = entry.evaluation
			}
		case UpperBound:
			if entry.evaluation < beta {
				beta = entry.evaluation
			}
		}
		if alpha >= beta {
			return entry.evaluation, true
		}
	}

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves)

	if len(moves) == 0 {
		if board.InCheck() {
			return -(Mate + Evaluation(depth)), true
		}
		return 0, true
	}

	alphaOrig := alpha

	if depth == 0 {
		return quiescence(board, alpha, beta, deadline)
	}

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, -beta, -alpha, deadline)
		if !ok {
			return 0, false
		}
		score = -score

		if score >= beta {
			Store(TableEntry{zobristHash: hash, depth: depth, evaluation: beta, flag: LowerBound, bestMove: move})
			return beta, true
		}
		if score > alpha {
			alpha = score
			bestMove = move
		}
	}

	flag := UpperBound
	if alpha > alphaOrig {
		flag = Exact
	}

	Store(TableEntry{zobristHash: hash, depth: depth, evaluation: alpha, flag: flag, bestMove: bestMove})
	return alpha, true
}
