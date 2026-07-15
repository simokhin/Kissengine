package main

import "testing"

func findMoveTo(moves []Move, to Square) (Move, bool) {
	for _, m := range moves {
		if m.To() == to {
			return m, true
		}
	}
	return 0, false
}

func TestGenerateKingMoves(t *testing.T) {
	tests := []struct {
		name             string
		inputFEN         string
		inputFrom        string
		expectedMovesLen int
		expectedMoves    []struct {
			to            string
			flag          Flag
			capturedPiece Piece
		}
	}{
		{
			"the king is on the edge of the board",
			"8/8/8/8/8/8/8/K7 w - - 0 1",
			"a1",
			3,
			nil,
		},
		{
			"the king is in the center of the board",
			"8/8/8/8/3K4/8/8/8 w - - 0 1",
			"d4",
			8,
			nil,
		},
		{
			"the king is blocked by its own piece",
			"8/8/8/8/8/8/P7/K7 w - - 0 1",
			"a1",
			2,
			nil,
		},
		{
			"the king can capture an enemy piece",
			"8/8/8/8/8/8/p7/K7 w - - 0 1",
			"a1",
			3,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"a2", Capture, Pawn | Black},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)
			from := FileRankToSquareIndex(SquareNotationToFileRank(tt.inputFrom))

			moves := GenerateJumpingPieceMoves(from, board, King)
			if len(moves) != tt.expectedMovesLen {
				t.Errorf("want %d moves, got %d moves", tt.expectedMovesLen, len(moves))
			}

			for _, expected := range tt.expectedMoves {
				expectedFile, expectedRank := SquareNotationToFileRank(expected.to)
				expectedTo := FileRankToSquareIndex(expectedFile, expectedRank)

				move, found := findMoveTo(moves, expectedTo)
				if !found {
					t.Errorf("expected a move to %s, but none was found", expected.to)
					continue
				}

				if move.Flag() != expected.flag {
					t.Errorf("move to %s: want flag %d, got %d", expected.to, expected.flag, move.Flag())
				}

				if move.CapturedPiece() != expected.capturedPiece {
					t.Errorf("move to %s: want captured piece %d, got %d", expected.to, expected.capturedPiece, move.CapturedPiece())
				}
			}
		})
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	tests := []struct {
		name             string
		inputFEN         string
		inputFrom        string
		expectedMovesLen int
		expectedMoves    []struct {
			to            string
			flag          Flag
			capturedPiece Piece
		}
	}{
		{
			"the knight is in the corner of the board",
			"8/8/8/8/8/8/8/N7 w - - 0 1",
			"a1",
			2,
			nil,
		},
		{
			"the knight is in the center of the board",
			"8/8/8/8/3N4/8/8/8 w - - 0 1",
			"d4",
			8,
			nil,
		},
		{
			"the knight is blocked by its own piece",
			"8/8/8/8/8/8/2N5/N7 w - - 0 1",
			"a1",
			1,
			nil,
		},
		{
			"the knight can capture an enemy piece",
			"8/8/8/8/8/8/2n5/N7 w - - 0 1",
			"a1",
			2,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"c2", Capture, Knight | Black},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)
			from := FileRankToSquareIndex(SquareNotationToFileRank(tt.inputFrom))

			moves := GenerateJumpingPieceMoves(from, board, Knight)
			if len(moves) != tt.expectedMovesLen {
				t.Errorf("want %d moves, got %d moves", tt.expectedMovesLen, len(moves))
			}

			for _, expected := range tt.expectedMoves {
				expectedFile, expectedRank := SquareNotationToFileRank(expected.to)
				expectedTo := FileRankToSquareIndex(expectedFile, expectedRank)

				move, found := findMoveTo(moves, expectedTo)
				if !found {
					t.Errorf("expected a move to %s, but none was found", expected.to)
					continue
				}

				if move.Flag() != expected.flag {
					t.Errorf("move to %s: want flag %d, got %d", expected.to, expected.flag, move.Flag())
				}

				if move.CapturedPiece() != expected.capturedPiece {
					t.Errorf("move to %s: want captured piece %d, got %d", expected.to, expected.capturedPiece, move.CapturedPiece())
				}
			}
		})
	}
}
