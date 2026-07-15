package main

var KingOffsets = [8]Square{-1, +1, -15, +15, -16, +16, -17, +17}

func GenerateKingMoves(from Square, board BoardState) []Move {
	offsets := KingOffsets

	moves := []Move{}

	kingsColor := board.squares[from].Color()

	for i := range offsets {
		candidate := from + offsets[i]

		if candidate&0x88 != 0 {
			continue
		}

		piece := board.squares[candidate]

		if piece == Empty {
			newMove := NewMove(from, candidate, 0, QuietMove, 0)
			moves = append(moves, newMove)
		} else {
			targetColor := piece.Color()

			if targetColor != kingsColor {
				newMove := NewMove(from, candidate, 0, Capture, piece)
				moves = append(moves, newMove)
			} else {
				continue
			}
		}
	}

	return moves
}
