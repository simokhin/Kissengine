package engine

import (
	"sort"
	"time"
)

const (
	Infinity      Evaluation = 1_000_000
	Mate                     = 100_000
	MateThreshold            = Mate - 1000
)

func moveScore(board BoardState, move Move, ttMove Move) int {
	if move == ttMove {
		return 1_000_000
	}
	if move.CapturedPiece() == Empty {
		return 0
	}
	attacker := board.PieceAt(move.From())
	return int(pieceValues[move.CapturedPiece().Type()])*10 - int(pieceValues[attacker.Type()])
}

func orderMoves(board BoardState, moves []Move, ttMove Move) []Move {
	sort.Slice(moves, func(i, j int) bool {
		return moveScore(board, moves[i], ttMove) > moveScore(board, moves[j], ttMove)
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
	var ply int

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves, Move(0))

	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, ply+1, -Infinity, -best, deadline)
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

func quiescence(board BoardState, ply int, alpha, beta Evaluation, deadline time.Time) (Evaluation, bool) {
	if !deadline.IsZero() && time.Now().After(deadline) {
		return 0, false
	}

	// quiescence has no search-depth parameter of its own, so entries it stores always use depth 0 —
	// that also makes any entry from negaMax (depth >= 0) trustworthy here, since a fuller search is only better.
	hash := ComputeHash(board)
	entry, found := Probe(hash)
	adjustedEval := entry.evaluation
	if adjustedEval > MateThreshold {
		adjustedEval -= Evaluation(ply)
	} else if adjustedEval < -MateThreshold {
		adjustedEval += Evaluation(ply)
	}
	if found {
		switch entry.flag {
		case Exact:
			return adjustedEval, true
		case LowerBound:
			if adjustedEval > alpha {
				alpha = adjustedEval
			}
		case UpperBound:
			if adjustedEval < beta {
				beta = adjustedEval
			}
		}
		if alpha >= beta {
			return adjustedEval, true
		}
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
	orderMoves(board, moves, entry.bestMove)

	if inCheck && len(moves) == 0 {
		return -Mate, true
	}

	alphaOrig := alpha
	var bestMove Move

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := quiescence(newBoard, ply+1, -beta, -alpha, deadline)
		if !ok {
			return 0, false
		}
		score = -score

		if score >= beta {
			storedEval := beta
			if storedEval > MateThreshold {
				storedEval += Evaluation(ply)
			} else if storedEval < -MateThreshold {
				storedEval -= Evaluation(ply)
			}
			Store(TableEntry{zobristHash: hash, depth: 0, evaluation: storedEval, flag: LowerBound, bestMove: move})
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
	storedEval := alpha
	if storedEval > MateThreshold {
		storedEval += Evaluation(ply)
	} else if storedEval < -MateThreshold {
		storedEval -= Evaluation(ply)
	}
	Store(TableEntry{zobristHash: hash, depth: 0, evaluation: storedEval, flag: flag, bestMove: bestMove})

	return alpha, true
}

func negaMax(board BoardState, depth int, ply int, alpha, beta Evaluation, deadline time.Time) (Evaluation, bool) {
	if !deadline.IsZero() && time.Now().After(deadline) {
		return 0, false
	}

	if depth == 0 && !board.InCheck() {
		return quiescence(board, ply, alpha, beta, deadline)
	}

	var bestMove Move

	hash := ComputeHash(board)
	entry, found := Probe(hash)
	adjustedEval := entry.evaluation
	if adjustedEval > MateThreshold {
		adjustedEval -= Evaluation(ply)
	} else if adjustedEval < -MateThreshold {
		adjustedEval += Evaluation(ply)
	}
	if found && entry.depth >= depth {
		switch entry.flag {
		case Exact:
			return adjustedEval, true
		case LowerBound:
			if adjustedEval > alpha {
				alpha = adjustedEval
			}
		case UpperBound:
			if adjustedEval < beta {
				beta = adjustedEval
			}
		}
		if alpha >= beta {
			return adjustedEval, true
		}
	}

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves, entry.bestMove)

	if len(moves) == 0 {
		if board.InCheck() {
			return -(Mate - Evaluation(ply)), true
		}
		return 0, true
	}

	alphaOrig := alpha

	if depth == 0 {
		return quiescence(board, ply, alpha, beta, deadline)
	}

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, ply+1, -beta, -alpha, deadline)
		if !ok {
			return 0, false
		}
		score = -score

		if score >= beta {
			storedEval := beta
			if storedEval > MateThreshold {
				storedEval += Evaluation(ply)
			} else if storedEval < -MateThreshold {
				storedEval -= Evaluation(ply)
			}
			Store(TableEntry{zobristHash: hash, depth: depth, evaluation: storedEval, flag: LowerBound, bestMove: move})
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

	storedEval := alpha
	if storedEval > MateThreshold {
		storedEval += Evaluation(ply)
	} else if storedEval < -MateThreshold {
		storedEval -= Evaluation(ply)
	}
	Store(TableEntry{zobristHash: hash, depth: depth, evaluation: storedEval, flag: flag, bestMove: bestMove})
	return alpha, true
}
