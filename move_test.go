package main

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
