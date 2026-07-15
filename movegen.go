package main

import "slices"

var KingOffsets = [8]Square{-1, +1, -15, +15, -16, +16, -17, +17}
var KnightOffsets = [8]Square{-14, +14, -18, +18, -31, +31, -33, +33}
var BishopOffsets = [4]Square{+15, -15, +17, -17}
var RookOffsets = [4]Square{+1, -1, +16, -16}
var WhitePawnOffsets = [4]Square{+16, +32, +15, +17}
var BlackPawnOffsets = [4]Square{-16, -32, -15, -17}

var piecesToPromote = [4]Piece{
	Bishop, Knight, Rook, Queen,
}

func GeneratePawnMoves(from Square, board BoardState) []Move {
	var offsets []Square
	var moves []Move
	var promotionRank int
	var startingRank int

	pieceColor := board.squares[from].Color()
	switch pieceColor {
	case White:
		offsets = WhitePawnOffsets[:]
		startingRank = Rank2
		promotionRank = Rank8
	case Black:
		offsets = BlackPawnOffsets[:]
		startingRank = Rank7
		promotionRank = Rank1
	}

	for i := range offsets {
		if i == 0 {
			candidateSquare := from + offsets[i]

			if candidateSquare&0x88 == 0 {
				targetPiece := board.squares[candidateSquare]
				if targetPiece != Empty {
					continue
				} else {
					_, rank := SquareIndexToFileRank(candidateSquare)
					if rank == promotionRank {
						for _, p := range piecesToPromote {
							newMove := NewMove(from, candidateSquare, p, QuietMove, Empty)
							moves = append(moves, newMove)
						}
					} else {
						newMove := NewMove(from, candidateSquare, 0, QuietMove, Empty)
						moves = append(moves, newMove)
					}
				}
			}
		} else if i == 1 {
			candidateSquare := from + offsets[i]

			intermediateSquare := from + offsets[0]
			if board.squares[intermediateSquare] != Empty {
				continue
			}

			_, rank := SquareIndexToFileRank(from)
			if rank != startingRank {
				continue
			} else {
				targetPiece := board.squares[candidateSquare]
				if targetPiece != Empty {
					continue
				} else {
					newMove := NewMove(from, candidateSquare, 0, DoublePawnMove, Empty)
					moves = append(moves, newMove)
				}
			}
		} else if i == 2 || i == 3 {
			candidateSquare := from + offsets[i]

			if candidateSquare&0x88 != 0 {
				continue
			}

			targetPiece := board.squares[candidateSquare]

			if targetPiece == Empty {
				if candidateSquare == board.enPassantSquare {
					capturedPawnSquare := candidateSquare - offsets[0]
					capturedPawn := board.squares[capturedPawnSquare]

					newMove := NewMove(from, candidateSquare, 0, EnPassantCapture, capturedPawn)
					moves = append(moves, newMove)
				}
				continue
			}

			if targetPiece.Color() == pieceColor {
				continue
			}

			_, rank := SquareIndexToFileRank(candidateSquare)
			if rank == promotionRank {
				for _, p := range piecesToPromote {
					newMove := NewMove(from, candidateSquare, p, Capture, targetPiece)
					moves = append(moves, newMove)
				}
			} else {
				newMove := NewMove(from, candidateSquare, 0, Capture, targetPiece)
				moves = append(moves, newMove)
			}
		}
	}

	return moves
}

func GenerateSlidingPieceMoves(from Square, board BoardState, piece Piece) []Move {
	var offsets []Square

	switch piece {
	case Bishop:
		offsets = BishopOffsets[:]
	case Rook:
		offsets = RookOffsets[:]
	case Queen:
		offsets = slices.Concat(BishopOffsets[:], RookOffsets[:])
	}

	moves := []Move{}

	pieceColor := board.squares[from].Color()
	for i := range offsets {
		candidateSquare := from
		for {
			candidateSquare += offsets[i]

			if candidateSquare&0x88 == 0 {
				targetPiece := board.squares[candidateSquare]

				if targetPiece == Empty {
					newMove := NewMove(from, candidateSquare, 0, QuietMove, 0)
					moves = append(moves, newMove)
				} else {
					if targetPiece.Color() != pieceColor {
						newMove := NewMove(from, candidateSquare, 0, Capture, targetPiece)
						moves = append(moves, newMove)
						break
					} else {
						break
					}
				}
			} else {
				break
			}
		}
	}

	return moves
}

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
