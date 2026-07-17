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
	Score Evaluation
}

const (
	Infinity      Evaluation = 1_000_000
	Mate                     = 100_000
	MateThreshold            = Mate - 1000
	killerBonus              = 50
)

var killerMoves [128][2]Move
var historyHeuristic [128][128]int

func moveScore(board BoardState, move Move, ttMove Move, killer1, killer2 Move) int {
	if move == ttMove {
		return 1_000_000
	}
	if move.CapturedPiece() == Empty {
		if move == killer1 || move == killer2 {
			return killerBonus
		}
		return historyHeuristic[move.From()][move.To()]
	}
	attacker := board.PieceAt(move.From())
	return int(pieceValues[move.CapturedPiece().Type()])*10 - int(pieceValues[attacker.Type()])
}

type scoredMove struct {
	move  Move
	score int
}

func orderMoves(board BoardState, moves []Move, ttMove Move, killer1, killer2 Move) []Move {
	scored := make([]scoredMove, len(moves))
	for i, move := range moves {
		scored[i] = scoredMove{move: move, score: moveScore(board, move, ttMove, killer1, killer2)}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	for i, s := range scored {
		moves[i] = s.move
	}

	return moves
}

func FindBestMove(board BoardState, depth int, history []ZobristHash) SearchResult {
	var searchResult SearchResult
	var nodes int
	move, score, _ := findBestMove(board, depth, time.Time{}, &nodes, history)

	searchResult.Move = move
	searchResult.Depth = depth
	searchResult.Nodes = nodes
	searchResult.Score = score

	return searchResult
}

func FindBestMoveByTime(board BoardState, timeLimit time.Duration, history []ZobristHash) SearchResult {
	deadline := time.Now().Add(timeLimit)
	var bestMove Move
	var bestDepth int
	var bestScore Evaluation
	var nodes int

	for depth := 1; ; depth++ {
		if time.Now().After(deadline) {
			break
		}
		move, score, ok := findBestMove(board, depth, deadline, &nodes, history)
		if !ok {
			break
		}
		bestMove = move
		bestDepth = depth
		bestScore = score
	}

	return SearchResult{Move: bestMove, Depth: bestDepth, Nodes: nodes, Score: bestScore}
}

func findBestMove(board BoardState, depth int, deadline time.Time, nodes *int, history []ZobristHash) (Move, Evaluation, bool) {
	var bestMove Move
	var ply int

	moves := GenerateLegalMoves(board)
	orderMoves(board, moves, Move(0), killerMoves[0][0], killerMoves[0][1])

	best := -Infinity

	for _, move := range moves {
		newBoard := MakeMove(board, move)
		score, ok := negaMax(newBoard, depth-1, ply+1, -Infinity, -best, deadline, nodes, history, true)
		if !ok {
			return bestMove, best, false
		}
		moveEvaluation := -score
		if moveEvaluation > best {
			best = moveEvaluation
			bestMove = move
		}
	}

	return bestMove, best, true
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
			if m.CapturedPiece() != Empty || m.Promotion() != Empty {
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

func hasNonPawnMaterial(board BoardState, color Piece) bool {
	for _, piece := range board.squares {
		if piece == Empty {
			continue
		}
		if piece.Color() == color && piece.Type() != Pawn && piece.Type() != King {
			return true
		}
	}
	return false
}

func negaMax(board BoardState, depth int, ply int, alpha, beta Evaluation, deadline time.Time, nodes *int, history []ZobristHash, allowNull bool) (Evaluation, bool) {
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

	if allowNull && depth >= 3 && !board.InCheck() && hasNonPawnMaterial(board, board.sideToMove.Color()) {
		nullBoard := board
		if nullBoard.sideToMove == WhiteToMove {
			nullBoard.sideToMove = BlackToMove
		} else {
			nullBoard.sideToMove = WhiteToMove
		}
		nullBoard.enPassantSquare = NoSquare

		R := 2
		if depth >= 6 {
			R = 3
		}
		score, ok := negaMax(nullBoard, depth-1-R, ply+1, -beta, -beta+1, deadline, nodes, history, false)
		if !ok {
			return 0, false
		}
		score = -score
		if score >= beta {
			// The null-move result alone is not trustworthy enough on its own — verify it with a
			// real (non-null), reduced-depth search of the actual position before relying on it.
			// Without this, unverified null-move cutoffs occasionally hide a zugzwang-like error
			// that the shallow null search missed, causing wildly inconsistent node counts between
			// neighboring depths.
			verifyScore, ok := negaMax(board, depth-R, ply, alpha, beta, deadline, nodes, history, false)
			if !ok {
				return 0, false
			}
			if verifyScore >= beta {
				return beta, true
			}
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
		score, ok := negaMax(newBoard, depth-1, ply+1, -beta, -alpha, deadline, nodes, childHistory, true)
		if !ok {
			return 0, false
		}
		score = -score

		if score >= beta {
			if move.CapturedPiece() == Empty {
				historyHeuristic[move.From()][move.To()] += depth * depth
				if move != killerMoves[ply][0] {
					killerMoves[ply][1] = killerMoves[ply][0]
					killerMoves[ply][0] = move
				}
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
