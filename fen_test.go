package main

import "testing"

func TestParseFEN(t *testing.T) {
	expectedPieces := []struct {
		notation string
		piece    int8
	}{
		{"a1", Rook | White},
		{"b1", Knight | White},
		{"c1", Bishop | White},
		{"d1", Queen | White},
		{"e1", King | White},
		{"f1", Bishop | White},
		{"g1", Knight | White},
		{"h1", Rook | White},

		{"a2", Pawn | White},
		{"b2", Pawn | White},
		{"c2", Pawn | White},
		{"d2", Pawn | White},
		{"e2", Pawn | White},
		{"f2", Pawn | White},
		{"g2", Pawn | White},
		{"h2", Pawn | White},

		{"a7", Pawn | Black},
		{"b7", Pawn | Black},
		{"c7", Pawn | Black},
		{"d7", Pawn | Black},
		{"e7", Pawn | Black},
		{"f7", Pawn | Black},
		{"g7", Pawn | Black},
		{"h7", Pawn | Black},

		{"a8", Rook | Black},
		{"b8", Knight | Black},
		{"c8", Bishop | Black},
		{"d8", Queen | Black},
		{"e8", King | Black},
		{"f8", Bishop | Black},
		{"g8", Knight | Black},
		{"h8", Rook | Black},
	}

	board := ParseFEN(StartFen)

	for _, expected := range expectedPieces {
		file, rank := SquareNotationToFileRank(expected.notation)
		squareIndex := FileRankToSquareIndex(file, rank)

		if board.squares[squareIndex] != expected.piece {
			t.Errorf("Expected piece to be %d but was %d", expected.piece, board.squares[squareIndex])
		}
	}

	if board.sideToMove != WhiteToMove {
		t.Fatalf("Expected side to move %d, got %d", WhiteToMove, board.sideToMove)
	}

	if board.castleRights != (WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide) {
		t.Fatalf("Expected all castle rights set, got %d", board.castleRights)
	}

	if board.enPassantSquare != NoSquare {
		t.Fatalf("Expected %d, got %d", NoSquare, board.enPassantSquare)
	}

	if board.fiftyMovesRuleCount != 0 {
		t.Fatalf("Expected %d, got %d", 0, board.fiftyMovesRuleCount)
	}

	if board.movesCount != 1 {
		t.Fatalf("Expected moves to be 1, got %d", board.movesCount)
	}
}
