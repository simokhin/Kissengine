package engine

import "math/bits"

var pieceValues = [7]Evaluation{
	Pawn:   100,
	Knight: 300,
	Bishop: 300,
	Rook:   500,
	Queen:  900,
	King:   0,
}

var kingAttackersValues = []Evaluation{
	Knight: 2,
	Bishop: 2,
	Rook:   3,
	Queen:  5,
}

const (
	bishopPairBonus          Evaluation = 30
	doubledPawnPenalty       Evaluation = 15
	isolatedPawnPenalty      Evaluation = 15
	openFileBonus            Evaluation = 20
	semiOpenFileBonus        Evaluation = 10
	missingShieldPawnPenalty Evaluation = 10
	mobilityBonus            Evaluation = 3
	kingAttackUnitValue      Evaluation = 8
)

var passedPawnBonus = [8]Evaluation{0, 5, 10, 20, 35, 60, 100, 0}

type Evaluation int

// aheadMask returns a bitmask of the ranks strictly ahead of rank, from color's point of view
// (higher ranks for White, lower ranks for Black).
func aheadMask(rank int, color Piece) uint8 {
	if color == White {
		return uint8(0xFF) << (rank + 1)
	}
	return uint8(0xFF) >> (8 - rank)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// mobility counts the squares a knight/bishop/rook/queen could move to right now — either
// empty squares or squares occupied by an enemy piece. It doesn't build real Move values
// (that would be needlessly expensive here), just counts reachable destinations.
//
// Along the way it also tallies kingAttackWeight: a weighted count of how many of those
// destinations land in the enemy king's zone (the king square plus its 8 neighbors), as a
// proxy for how much direct piece pressure the enemy king is under right now.
func mobility(board BoardState, from Square, piece Piece) (count int, kingAttackWeight Evaluation) {
	pieceColor := piece.Color()
	var enemyKingFile, enemyKingRank int

	if pieceColor == White {
		enemyKingFile, enemyKingRank = SquareIndexToFileRank(board.blackKingSquare)
	} else {
		enemyKingFile, enemyKingRank = SquareIndexToFileRank(board.whiteKingSquare)
	}

	inEnemyKingZone := func(to Square) bool {
		toFile, toRank := SquareIndexToFileRank(to)
		return abs(toFile-enemyKingFile) <= 1 && abs(toRank-enemyKingRank) <= 1
	}

	switch piece.Type() {
	case Knight:
		for _, offset := range KnightOffsets {
			to := from + offset
			if !to.IsOnBoard() {
				continue
			}
			target := board.squares[to]
			if target == Empty || target.Color() != pieceColor {
				count++
				if inEnemyKingZone(to) {
					kingAttackWeight += kingAttackersValues[piece.Type()]
				}
			}
		}
	case Bishop, Rook, Queen:
		var offsets []Square
		switch piece.Type() {
		case Bishop:
			offsets = BishopOffsets[:]
		case Rook:
			offsets = RookOffsets[:]
		case Queen:
			offsets = QueenOffsets[:]
		}
		for _, offset := range offsets {
			to := from
			for {
				to += offset
				if !to.IsOnBoard() {
					break
				}
				target := board.squares[to]
				if target == Empty {
					count++
					if inEnemyKingZone(to) {
						kingAttackWeight += kingAttackersValues[piece.Type()]
					}
					continue
				}
				if target.Color() != pieceColor {
					count++
					if inEnemyKingZone(to) {
						kingAttackWeight += kingAttackersValues[piece.Type()]
					}
				}
				break
			}
		}
	}

	return count, kingAttackWeight
}

// kingShieldPenalty checks the king's own file and the two adjacent files for a friendly
// pawn somewhere ahead of the king (using the same "ahead" notion as passed pawns) — one
// missing shield file adds a fixed penalty.
func kingShieldPenalty(kingSquare Square, color Piece, ownRanks [8]uint8) Evaluation {
	file, rank := SquareIndexToFileRank(kingSquare)
	mask := aheadMask(rank, color)

	var penalty Evaluation
	for _, f := range [3]int{file - 1, file, file + 1} {
		if f < 0 || f > 7 {
			continue
		}
		if ownRanks[f]&mask == 0 {
			penalty += missingShieldPawnPenalty
		}
	}
	return penalty
}

func Evaluate(board BoardState) Evaluation {
	var evaluation Evaluation
	phase := gamePhase(board)
	var whiteBishops, blackBishops int
	var whitePawnRanks, blackPawnRanks [8]uint8
	var pawnSquares []Square
	var rookSquares []Square
	var whiteKingAttackWeight, blackKingAttackWeight Evaluation

	for i := range board.squares {
		piece := board.squares[i]
		if piece == Empty {
			continue
		}

		if piece.Type() == Pawn {
			file, rank := SquareIndexToFileRank(Square(i))
			if piece.Color() == White {
				whitePawnRanks[file] |= 1 << rank
			} else {
				blackPawnRanks[file] |= 1 << rank
			}
			pawnSquares = append(pawnSquares, Square(i))
		}

		if piece.Type() == Bishop {
			if piece.Color() == White {
				whiteBishops++
			} else {
				blackBishops++
			}
		}

		if piece.Type() == Rook {
			rookSquares = append(rookSquares, Square(i))
		}

		if piece.Type() == Knight || piece.Type() == Bishop || piece.Type() == Rook || piece.Type() == Queen {
			moveCount, kingAttackWeight := mobility(board, Square(i), piece)
			mobilityScore := Evaluation(moveCount) * mobilityBonus
			if piece.Color() == board.SideToMove().Color() {
				evaluation += mobilityScore
			} else {
				evaluation -= mobilityScore
			}

			if piece.Color() == White {
				blackKingAttackWeight += kingAttackWeight
			} else {
				whiteKingAttackWeight += kingAttackWeight
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

	for _, square := range pawnSquares {
		piece := board.squares[square]
		file, rank := SquareIndexToFileRank(square)

		ownRanks, opponentRanks := &whitePawnRanks, &blackPawnRanks
		if piece.Color() == Black {
			ownRanks, opponentRanks = &blackPawnRanks, &whitePawnRanks
		}

		var bonus Evaluation

		if bits.OnesCount8(ownRanks[file]) > 1 {
			bonus -= doubledPawnPenalty
		}

		isolated := true
		if file > 0 && ownRanks[file-1] != 0 {
			isolated = false
		}
		if file < 7 && ownRanks[file+1] != 0 {
			isolated = false
		}
		if isolated {
			bonus -= isolatedPawnPenalty
		}

		mask := aheadMask(rank, piece.Color())
		passed := true
		for _, f := range [3]int{file - 1, file, file + 1} {
			if f < 0 || f > 7 {
				continue
			}
			if opponentRanks[f]&mask != 0 {
				passed = false
				break
			}
		}
		if passed {
			promotionRank := rank
			if piece.Color() == Black {
				promotionRank = 7 - rank
			}
			bonus += passedPawnBonus[promotionRank]
		}

		if piece.Color() == board.SideToMove().Color() {
			evaluation += bonus
		} else {
			evaluation -= bonus
		}
	}

	for _, square := range rookSquares {
		piece := board.squares[square]
		file, _ := SquareIndexToFileRank(square)

		var bonus Evaluation
		if whitePawnRanks[file] == 0 && blackPawnRanks[file] == 0 {
			bonus = openFileBonus
		} else {
			ownRanks := whitePawnRanks
			if piece.Color() == Black {
				ownRanks = blackPawnRanks
			}
			if ownRanks[file] == 0 {
				bonus = semiOpenFileBonus
			}
		}

		if piece.Color() == board.SideToMove().Color() {
			evaluation += bonus
		} else {
			evaluation -= bonus
		}
	}

	whiteShieldPenalty := Evaluation(phase * float64(kingShieldPenalty(board.whiteKingSquare, White, whitePawnRanks)))
	blackShieldPenalty := Evaluation(phase * float64(kingShieldPenalty(board.blackKingSquare, Black, blackPawnRanks)))

	whiteKingAttackPenalty := Evaluation(phase * float64(whiteKingAttackWeight*kingAttackUnitValue))
	blackKingAttackPenalty := Evaluation(phase * float64(blackKingAttackWeight*kingAttackUnitValue))

	if board.SideToMove().Color() == White {
		evaluation -= whiteShieldPenalty
		evaluation += blackShieldPenalty
		evaluation -= whiteKingAttackPenalty
		evaluation += blackKingAttackPenalty
	} else {
		evaluation -= blackShieldPenalty
		evaluation += whiteShieldPenalty
		evaluation -= blackKingAttackPenalty
		evaluation += whiteKingAttackPenalty
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
