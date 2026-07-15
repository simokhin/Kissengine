package main

var KingOffsets = [8]Square{-1, +1, -15, +15, -16, +16, -17, +17}
var KnightOffsets = [8]Square{-14, +14, -18, +18, -31, +31, -33, +33}

func GenerateJumpingPieceMoves(from Square, board BoardState, piece Piece) []Move {
	var offsets []Square
	switch piece {
	case King:
		offsets = KingOffsets[:]
	case Knight:
		offsets = KnightOffsets[:]
	}

	moves := []Move{}

	pieceColor := board.squares[from].Color()

	for i := range offsets {
		candidateSquare := from + offsets[i]

		if candidateSquare&0x88 != 0 {
			continue
		}

		targetPiece := board.squares[candidateSquare]

		if targetPiece == Empty {
			newMove := NewMove(from, candidateSquare, 0, QuietMove, 0)
			moves = append(moves, newMove)
		} else {
			targetColor := targetPiece.Color()

			if targetColor != pieceColor {
				newMove := NewMove(from, candidateSquare, 0, Capture, targetPiece)
				moves = append(moves, newMove)
			} else {
				continue
			}
		}
	}

	return moves

}

