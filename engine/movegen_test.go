package engine

import "testing"

func findMoveTo(moves []Move, to Square) (Move, bool) {
	for _, m := range moves {
		if m.To() == to {
			return m, true
		}
	}
	return 0, false
}

func findMoveToWithFlag(moves []Move, to Square, flag Flag) (Move, bool) {
	for _, m := range moves {
		if m.To() == to && m.Flag() == flag {
			return m, true
		}
	}
	return 0, false
}

func findMove(moves []Move, from, to Square, flag Flag) (Move, bool) {
	for _, m := range moves {
		if m.From() == from && m.To() == to && m.Flag() == flag {
			return m, true
		}
	}
	return 0, false
}

func TestGenerateLegalMoves(t *testing.T) {
	type moveExpectation struct {
		from, to string
		flag     Flag
	}

	tests := []struct {
		name             string
		inputFEN         string
		expectedMovesLen int
		mustInclude      []moveExpectation
		mustExclude      []moveExpectation
	}{
		{
			"starting position: no king is in danger, all pseudo-legal moves stay legal",
			StartFen,
			20,
			[]moveExpectation{
				{"e2", "e4", DoublePawnMove},
				{"b1", "c3", QuietMove},
			},
			nil,
		},
		{
			"king cannot move onto or through a square attacked by an enemy rook",
			"8/8/8/8/8/8/5r2/4K3 w - - 0 1",
			2,
			[]moveExpectation{
				{"e1", "d1", QuietMove},
				{"e1", "f2", Capture},
			},
			[]moveExpectation{
				{"e1", "e2", QuietMove},
				{"e1", "f1", QuietMove},
				{"e1", "d2", QuietMove},
			},
		},
		{
			"when the king is in check, only moves that escape the check are legal",
			"4r3/8/8/8/8/8/8/R3K3 w - - 0 1",
			4,
			[]moveExpectation{
				{"e1", "d1", QuietMove},
				{"e1", "f1", QuietMove},
				{"e1", "d2", QuietMove},
				{"e1", "f2", QuietMove},
			},
			[]moveExpectation{
				{"e1", "e2", QuietMove},
				{"a1", "d1", QuietMove},
				{"a1", "b1", QuietMove},
				{"a1", "a4", QuietMove},
			},
		},
		{
			"blocking the check is legal, but unrelated moves of the same piece are not",
			"4r3/8/8/8/8/8/1B6/4K3 w - - 0 1",
			5,
			[]moveExpectation{
				{"e1", "d1", QuietMove},
				{"e1", "f1", QuietMove},
				{"e1", "d2", QuietMove},
				{"e1", "f2", QuietMove},
				{"b2", "e5", QuietMove},
			},
			[]moveExpectation{
				{"e1", "e2", QuietMove},
				{"b2", "c1", QuietMove},
				{"b2", "h8", QuietMove},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)

			moves := GenerateLegalMoves(board)
			if len(moves) != tt.expectedMovesLen {
				t.Errorf("want %d moves, got %d moves", tt.expectedMovesLen, len(moves))
			}

			for _, expected := range tt.mustInclude {
				fromFile, fromRank := SquareNotationToFileRank(expected.from)
				from := FileRankToSquareIndex(fromFile, fromRank)
				toFile, toRank := SquareNotationToFileRank(expected.to)
				to := FileRankToSquareIndex(toFile, toRank)

				if _, found := findMove(moves, from, to, expected.flag); !found {
					t.Errorf("expected a legal move %s->%s with flag %d, but none was found", expected.from, expected.to, expected.flag)
				}
			}

			for _, unexpected := range tt.mustExclude {
				fromFile, fromRank := SquareNotationToFileRank(unexpected.from)
				from := FileRankToSquareIndex(fromFile, fromRank)
				toFile, toRank := SquareNotationToFileRank(unexpected.to)
				to := FileRankToSquareIndex(toFile, toRank)

				if _, found := findMove(moves, from, to, unexpected.flag); found {
					t.Errorf("did not expect a legal move %s->%s with flag %d, but it was found", unexpected.from, unexpected.to, unexpected.flag)
				}
			}
		})
	}
}

func TestGenerateAllPseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name             string
		inputFEN         string
		expectedMovesLen int
		expectedMoves    []struct {
			to   string
			flag Flag
		}
	}{
		{
			"starting position generates all 20 opening moves for white",
			StartFen,
			20,
			[]struct {
				to   string
				flag Flag
			}{
				{"e4", DoublePawnMove},
				{"c3", QuietMove},
			},
		},
		{
			"starting position generates all 20 opening moves for black",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			20,
			[]struct {
				to   string
				flag Flag
			}{
				{"e5", DoublePawnMove},
				{"c6", QuietMove},
			},
		},
		{
			"only pieces belonging to the side to move generate moves",
			"8/8/8/8/8/8/8/N6n w - - 0 1",
			2,
			nil,
		},
		{
			"combines rook, king, and castling moves for a mixed position",
			"8/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			26,
			[]struct {
				to   string
				flag Flag
			}{
				{"g1", KingCastle},
				{"c1", QueenCastle},
				{"a8", QuietMove},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)

			moves := GenerateAllPseudoLegalMoves(board)
			if len(moves) != tt.expectedMovesLen {
				t.Errorf("want %d moves, got %d moves", tt.expectedMovesLen, len(moves))
			}

			for _, expected := range tt.expectedMoves {
				expectedFile, expectedRank := SquareNotationToFileRank(expected.to)
				expectedTo := FileRankToSquareIndex(expectedFile, expectedRank)

				_, found := findMoveToWithFlag(moves, expectedTo, expected.flag)
				if !found {
					t.Errorf("expected a move to %s with flag %d, but none was found", expected.to, expected.flag)
				}
			}
		})
	}
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

			moves := GenerateJumpingPieceMoves(from, board, King, nil)
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

			moves := GenerateJumpingPieceMoves(from, board, Knight, nil)
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

func TestGenerateCastlingMoves(t *testing.T) {
	tests := []struct {
		name             string
		inputFEN         string
		expectedMovesLen int
		expectedMoves    []struct {
			to   string
			flag Flag
		}
	}{
		{
			"white king-side castle available when the path is clear and safe",
			"8/8/8/8/8/8/8/4K2R w K - 0 1",
			1,
			[]struct {
				to   string
				flag Flag
			}{
				{"g1", KingCastle},
			},
		},
		{
			"white king-side castle blocked by a piece between king and rook",
			"8/8/8/8/8/8/8/4KN1R w K - 0 1",
			0,
			nil,
		},
		{
			"white king-side castle blocked when the king is in check",
			"4r3/8/8/8/8/8/8/4K2R w K - 0 1",
			0,
			nil,
		},
		{
			"white king-side castle blocked when the king passes through an attacked square",
			"6r1/8/8/8/8/8/8/4K2R w K - 0 1",
			0,
			nil,
		},
		{
			"white queen-side castle available when the path is clear and safe",
			"8/8/8/8/8/8/8/R3K3 w Q - 0 1",
			1,
			[]struct {
				to   string
				flag Flag
			}{
				{"c1", QueenCastle},
			},
		},
		{
			"white queen-side castle blocked by a piece on b1",
			"8/8/8/8/8/8/8/RN2K3 w Q - 0 1",
			0,
			nil,
		},
		{
			"no castling moves when castle rights are not set",
			"8/8/8/8/8/8/8/4K2R w - - 0 1",
			0,
			nil,
		},
		{
			"black king-side and queen-side castle available",
			"r3k2r/8/8/8/8/8/8/8 b kq - 0 1",
			2,
			[]struct {
				to   string
				flag Flag
			}{
				{"g8", KingCastle},
				{"c8", QueenCastle},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)

			moves := GenerateCastlingMoves(board, nil)
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

			moves := GenerateSlidingPieceMoves(from, board, tt.inputPiece, nil)
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

			moves := GeneratePawnMoves(from, board, nil)
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
