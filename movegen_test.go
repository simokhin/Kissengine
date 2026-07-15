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

func TestGenerateSlidingPieceMoves(t *testing.T) {
	tests := []struct {
		name             string
		inputFEN         string
		inputFrom        string
		inputPiece       Piece
		expectedMovesLen int
		expectedMoves    []struct {
			to            string
			flag          Flag
			capturedPiece Piece
		}
	}{
		{
			"the bishop is on an open board",
			"8/8/8/8/3B4/8/8/8 w - - 0 1",
			"d4",
			Bishop,
			13,
			nil,
		},
		{
			"the rook is on an open board",
			"8/8/8/8/3R4/8/8/8 w - - 0 1",
			"d4",
			Rook,
			14,
			nil,
		},
		{
			"the queen is on an open board",
			"8/8/8/8/3Q4/8/8/8 w - - 0 1",
			"d4",
			Queen,
			27,
			nil,
		},
		{
			"the rook is blocked by its own piece along a file",
			"8/8/8/8/P7/8/8/R7 w - - 0 1",
			"a1",
			Rook,
			9,
			nil,
		},
		{
			"the rook can capture along a file and stops there",
			"8/8/8/8/p7/8/8/R7 w - - 0 1",
			"a1",
			Rook,
			10,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"a4", Capture, Pawn | Black},
			},
		},
		{
			"the bishop slides several squares and stops after a capture",
			"8/8/8/8/3n4/8/8/B7 w - - 0 1",
			"a1",
			Bishop,
			3,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"d4", Capture, Knight | Black},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)
			from := FileRankToSquareIndex(SquareNotationToFileRank(tt.inputFrom))

			moves := GenerateSlidingPieceMoves(from, board, tt.inputPiece)
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

func TestGeneratePawnMoves(t *testing.T) {
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
			"white pawn single push on an open board",
			"8/8/8/8/8/4P3/8/8 w - - 0 1",
			"e3",
			1,
			nil,
		},
		{
			"white pawn double push from the starting rank",
			"8/8/8/8/8/8/4P3/8 w - - 0 1",
			"e2",
			2,
			nil,
		},
		{
			"white pawn is blocked and cannot push or double push",
			"8/8/8/8/8/4p3/4P3/8 w - - 0 1",
			"e2",
			0,
			nil,
		},
		{
			"white pawn can capture diagonally",
			"8/8/8/3n4/4P3/8/8/8 w - - 0 1",
			"e4",
			2,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"d5", Capture, Knight | Black},
			},
		},
		{
			"white pawn promotes by pushing and by capturing",
			"3n4/4P3/8/8/8/8/8/8 w - - 0 1",
			"e7",
			8,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"e8", QuietMove, Empty},
				{"d8", Capture, Knight | Black},
			},
		},
		{
			"white pawn can capture en passant",
			"8/8/8/3pP3/8/8/8/8 w - d6 0 1",
			"e5",
			2,
			[]struct {
				to            string
				flag          Flag
				capturedPiece Piece
			}{
				{"d6", EnPassantCapture, Pawn | Black},
			},
		},
		{
			"black pawn single and double push from the starting rank",
			"8/4p3/8/8/8/8/8/8 b - - 0 1",
			"e7",
			2,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)
			from := FileRankToSquareIndex(SquareNotationToFileRank(tt.inputFrom))

			moves := GeneratePawnMoves(from, board)
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
