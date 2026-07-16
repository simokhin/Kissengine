package engine

import "testing"

// TestMoveRoundTrip creates a new move and then extracts information from the created move
func TestMoveRoundTrip(t *testing.T) {
	from := FileRankToSquareIndex(SquareNotationToFileRank("e7"))
	to := FileRankToSquareIndex(SquareNotationToFileRank("e8"))

	move := NewMove(from, to, Queen, Capture, Rook|Black)

	gotTo := move.To()
	if gotTo != to {
		t.Fatalf("want %d, got %d", to, gotTo)
	}

	gotFrom := move.From()
	if gotFrom != from {
		t.Fatalf("want %d, got %d", from, gotFrom)
	}

	gotFlag := move.Flag()
	if gotFlag != Capture {
		t.Fatalf("want %d, got %d", Capture, gotFlag)
	}

	gotPromotion := move.Promotion()
	if gotPromotion != Queen {
		t.Fatalf("want %d, got %d", Queen, gotPromotion)
	}

	gotCapturedPiece := move.CapturedPiece()
	if gotCapturedPiece != Rook|Black {
		t.Fatalf("want %d, got %d", Rook|Black, gotCapturedPiece)
	}
}

func TestMakeMove(t *testing.T) {
	type squareExpectation struct {
		square string
		piece  Piece
	}

	tests := []struct {
		name                        string
		inputFEN                    string
		from, to                    string
		promotion                   Piece
		flag                        Flag
		capturedPiece               Piece
		expectedSideToMove          SideToMove
		expectedCastleRights        CastleRights
		expectedEnPassantSquare     Square
		expectedFiftyMovesRuleCount int
		expectedMovesCount          int
		expectedSquares             []squareExpectation
	}{
		{
			"quiet knight move by white increments fifty-move count and keeps moves count",
			"8/8/8/8/8/8/8/N7 w - - 5 10",
			"a1", "b3",
			Empty, QuietMove, Empty,
			BlackToMove, 0, NoSquare, 6, 10,
			[]squareExpectation{
				{"a1", Empty},
				{"b3", White | Knight},
			},
		},
		{
			"quiet knight move by black increments moves count",
			"8/8/8/8/8/8/8/n7 b - - 3 7",
			"a1", "b3",
			Empty, QuietMove, Empty,
			WhiteToMove, 0, NoSquare, 4, 8,
			[]squareExpectation{
				{"a1", Empty},
				{"b3", Black | Knight},
			},
		},
		{
			"pawn push resets fifty-move count",
			"8/8/8/8/8/4P3/8/8 w - - 10 1",
			"e3", "e4",
			Empty, QuietMove, Empty,
			BlackToMove, 0, NoSquare, 0, 1,
			[]squareExpectation{
				{"e3", Empty},
				{"e4", White | Pawn},
			},
		},
		{
			"capture resets fifty-move count",
			"8/8/8/8/8/8/2p5/N7 w - - 8 4",
			"a1", "c2",
			Empty, Capture, Pawn | Black,
			BlackToMove, 0, NoSquare, 0, 4,
			[]squareExpectation{
				{"a1", Empty},
				{"c2", White | Knight},
			},
		},
		{
			"promotion places the promoted piece with the correct color",
			"8/4P3/8/8/8/8/8/8 w - - 0 1",
			"e7", "e8",
			Queen, QuietMove, Empty,
			BlackToMove, 0, NoSquare, 0, 1,
			[]squareExpectation{
				{"e7", Empty},
				{"e8", White | Queen},
			},
		},
		{
			"en passant capture removes the captured pawn",
			"8/8/8/3pP3/8/8/8/8 w - d6 0 1",
			"e5", "d6",
			Empty, EnPassantCapture, Pawn | Black,
			BlackToMove, 0, NoSquare, 0, 1,
			[]squareExpectation{
				{"e5", Empty},
				{"d6", White | Pawn},
				{"d5", Empty},
			},
		},
		{
			"double pawn push sets the en passant square for white",
			"8/8/8/8/8/8/4P3/8 w - - 0 1",
			"e2", "e4",
			Empty, DoublePawnMove, Empty,
			BlackToMove, 0, FileRankToSquareIndex(SquareNotationToFileRank("e3")), 0, 1,
			[]squareExpectation{
				{"e2", Empty},
				{"e4", White | Pawn},
			},
		},
		{
			"double pawn push sets the en passant square for black",
			"8/4p3/8/8/8/8/8/8 b - - 0 5",
			"e7", "e5",
			Empty, DoublePawnMove, Empty,
			WhiteToMove, 0, FileRankToSquareIndex(SquareNotationToFileRank("e6")), 0, 6,
			[]squareExpectation{
				{"e7", Empty},
				{"e5", Black | Pawn},
			},
		},
		{
			"white king-side castle moves the rook and clears white castle rights",
			"8/8/8/8/8/8/8/4K2R w KQ - 0 1",
			"e1", "g1",
			Empty, KingCastle, Empty,
			BlackToMove, 0, NoSquare, 1, 1,
			[]squareExpectation{
				{"e1", Empty},
				{"g1", White | King},
				{"h1", Empty},
				{"f1", White | Rook},
			},
		},
		{
			"black queen-side castle moves the rook and clears black castle rights",
			"r3k3/8/8/8/8/8/8/8 b kq - 0 1",
			"e8", "c8",
			Empty, QueenCastle, Empty,
			WhiteToMove, 0, NoSquare, 1, 2,
			[]squareExpectation{
				{"e8", Empty},
				{"c8", Black | King},
				{"a8", Empty},
				{"d8", Black | Rook},
			},
		},
		{
			"moving the king without castling clears both rights for that side",
			"8/8/8/8/8/8/8/4K3 w KQ - 0 1",
			"e1", "e2",
			Empty, QuietMove, Empty,
			BlackToMove, 0, NoSquare, 1, 1,
			[]squareExpectation{
				{"e1", Empty},
				{"e2", White | King},
			},
		},
		{
			"moving a rook from its corner clears only that right",
			"8/8/8/8/8/8/8/R3K3 w KQ - 0 1",
			"a1", "a4",
			Empty, QuietMove, Empty,
			BlackToMove, WhiteKingSide, NoSquare, 1, 1,
			[]squareExpectation{
				{"a1", Empty},
				{"a4", White | Rook},
			},
		},
		{
			"capturing an untouched rook in the corner clears that right",
			"4k2r/5N2/8/8/8/8/8/4K3 w KQkq - 0 1",
			"f7", "h8",
			Empty, Capture, Rook | Black,
			BlackToMove, WhiteKingSide | WhiteQueenSide | BlackQueenSide, NoSquare, 0, 1,
			[]squareExpectation{
				{"f7", Empty},
				{"h8", White | Knight},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)
			from := FileRankToSquareIndex(SquareNotationToFileRank(tt.from))
			to := FileRankToSquareIndex(SquareNotationToFileRank(tt.to))

			move := NewMove(from, to, tt.promotion, tt.flag, tt.capturedPiece)
			newBoard := MakeMove(board, move)

			if newBoard.sideToMove != tt.expectedSideToMove {
				t.Errorf("sideToMove: want %d, got %d", tt.expectedSideToMove, newBoard.sideToMove)
			}

			if newBoard.castleRights != tt.expectedCastleRights {
				t.Errorf("castleRights: want %d, got %d", tt.expectedCastleRights, newBoard.castleRights)
			}

			if newBoard.enPassantSquare != tt.expectedEnPassantSquare {
				t.Errorf("enPassantSquare: want %d, got %d", tt.expectedEnPassantSquare, newBoard.enPassantSquare)
			}

			if newBoard.fiftyMovesRuleCount != tt.expectedFiftyMovesRuleCount {
				t.Errorf("fiftyMovesRuleCount: want %d, got %d", tt.expectedFiftyMovesRuleCount, newBoard.fiftyMovesRuleCount)
			}

			if newBoard.movesCount != tt.expectedMovesCount {
				t.Errorf("movesCount: want %d, got %d", tt.expectedMovesCount, newBoard.movesCount)
			}

			for _, sq := range tt.expectedSquares {
				file, rank := SquareNotationToFileRank(sq.square)
				got := newBoard.squares[FileRankToSquareIndex(file, rank)]
				if got != sq.piece {
					t.Errorf("square %s: want piece %d, got %d", sq.square, sq.piece, got)
				}
			}
		})
	}
}
