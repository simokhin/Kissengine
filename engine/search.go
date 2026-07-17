package engine

import (
	"slices"
	"sort"
	"time"
)

type SearchResult struct {
	Move  Move
	Nodes int
	Depth int
}

const (
	Infinity      Evaluation = 1_000_000
	Mate                     = 100_000
	MateThreshold            = Mate - 1000
	killerBonus              = 50
)

var killerMoves [128][2]Move

func moveScore(board BoardState, move Move, ttMove Move, killer1, killer2 Move) int {
	if move == ttMove {
		return 1_000_000
	}
	if move.CapturedPiece() == Empty {
		if move == killer1 || move == killer2 {
			return killerBonus
		}
		return 0
	}
	attacker := board.PieceAt(move.From())
	return int(pieceValues[move.CapturedPiece().Type()])*10 - int(pieceValues[attacker.Type()])
}

func orderMoves(board BoardState, moves []Move, ttMove Move, killer1, killer2 Move) []Move {
	sort.Slice(moves, func(i, j int) bool {
		return moveScore(board, moves[i], ttMove, killer1, killer2) > moveScore(board, moves[j], ttMove, killer1, killer2)
	})
	return moves
}

func FindBestMove(board BoardState, depth int, history []ZobristHash) SearchResult {
	var searchResult SearchResult
	var nodes int
	move, _ := findBestMove(board, depth, time.Time{}, &nodes, history)

	searchResult.Move = move
	searchResult.Depth = depth
	searchResult.Nodes = nodes

	return searchResult
}

func FindBestMoveByTime(board BoardState, timeLimit time.Duration, history []ZobristHash) SearchResult {
	deadline := time.Now().Add(timeLimit)
	var bestMove Move
	var bestDepth int
	var nodes int

	for depth := 1; ; depth++ {
		if time.Now().After(deadline) {
			break
		}
		move, ok := findBestMove(board, depth, deadline, &nodes, history)
		if !ok {
			break
		}
		bestMove = move
		bestDepth = depth
	}

	return SearchResult{Move: bestMove, Depth: bestDepth, Nodes: nodes}
}

func findBestMove(board BoardState, depth int, deadline time.Time, nodes *int, history []ZobristHash) (Move, bool) {
	var bestMove Move
	var ply int

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves, Move(0), killerMoves[0][0], killerMoves[0][1])

	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, ply+1, -Infinity, -best, deadline, nodes, history)
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

func quiescence(board BoardState, ply int, alpha, beta Evaluation, deadline time.Time, nodes *int) (Evaluation, bool) {
	*nodes++

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
	orderMoves(board, moves, entry.bestMove, Move(0), Move(0))

	if inCheck && len(moves) == 0 {
		return -Mate, true
	}

	alphaOrig := alpha
	var bestMove Move

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := quiescence(newBoard, ply+1, -beta, -alpha, deadline, nodes)
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

func negaMax(board BoardState, depth int, ply int, alpha, beta Evaluation, deadline time.Time, nodes *int, history []ZobristHash) (Evaluation, bool) {
	*nodes++

	if !deadline.IsZero() && time.Now().After(deadline) {
		return 0, false
	}

	hash := ComputeHash(board)
	if slices.Contains(history, hash) {
		return 0, true
	}

	if depth == 0 && !board.InCheck() {
		return quiescence(board, ply, alpha, beta, deadline, nodes)
	}

	var bestMove Move

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
	orderMoves(board, moves, entry.bestMove, killerMoves[ply][0], killerMoves[ply][1])

	if len(moves) == 0 {
		if board.InCheck() {
			return -(Mate - Evaluation(ply)), true
		}
		return 0, true
	}

	alphaOrig := alpha

	if depth == 0 {
		return quiescence(board, ply, alpha, beta, deadline, nodes)
	}

	childHistory := append(history[:len(history):len(history)], hash)

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, ply+1, -beta, -alpha, deadline, nodes, childHistory)
		if !ok {
			return 0, false
		}
		score = -score

		if score >= beta {
			if move.CapturedPiece() == Empty && move != killerMoves[ply][0] {
				killerMoves[ply][1] = killerMoves[ply][0]
				killerMoves[ply][0] = move
			}

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
