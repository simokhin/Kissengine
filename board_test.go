package main

import "testing"

func TestFileRankToSquareIndex(t *testing.T) {
	tests := []struct {
		name                string
		inputFile           int
		inputRank           int
		expectedSquareIndex Square
	}{
		{"a1", FileA, Rank1, 0},
		{"h8", FileH, Rank8, 119},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			squareIndex := FileRankToSquareIndex(tt.inputFile, tt.inputRank)
			if squareIndex != tt.expectedSquareIndex {
				t.Fatalf("FileRankToSquareIndex(%d, %d) = %d; want %d", tt.inputFile, tt.inputRank, squareIndex, tt.expectedSquareIndex)
			}
		})
	}
}

func TestSquareIndexToFileRank(t *testing.T) {
	tests := []struct {
		name         string
		input        Square
		expectedFile int
		expectedRank int
	}{
		{"a1", 0, 0, 0},
		{"h8", 119, 7, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, rank := SquareIndexToFileRank(tt.input)
			if file != tt.expectedFile || rank != tt.expectedRank {
				t.Fatalf("IndexToFileRank(%d) = (%d, %d); want (%d, %d)", tt.input, file, rank, tt.expectedFile, tt.expectedRank)
			}
		})
	}
}

func TestFileRankToNotation(t *testing.T) {
	tests := []struct {
		name             string
		inputFile        int
		inputRank        int
		expectedNotation string
	}{
		{"a1", FileA, Rank1, "a1"},
		{"h8", FileH, Rank8, "h8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			squareNotation := FileRankToNotation(tt.inputFile, tt.inputRank)
			if squareNotation != tt.expectedNotation {
				t.Fatalf("FileRankToNotation(%d, %d) = %s; want %s", tt.inputFile, tt.inputRank, squareNotation, tt.expectedNotation)
			}
		})
	}
}

func TestSquareNotationToFileRank(t *testing.T) {
	tests := []struct {
		name                string
		inputSquareNotation string
		expectedFile        int
		expectedRank        int
	}{
		{"a1", "a1", FileA, Rank1},
		{"h8", "h8", FileH, Rank8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, rank := SquareNotationToFileRank(tt.inputSquareNotation)
			if file != tt.expectedFile || rank != tt.expectedRank {
				t.Fatalf("SquareNotationToFileRank(%s) = (%d, %d); want (%d, %d)", tt.inputSquareNotation, file, rank, tt.expectedFile, tt.expectedRank)
			}
		})
	}
}

func TestFileRankSquareIndexRoundTrip(t *testing.T) {
	for file := 0; file <= FileH; file++ {
		for rank := 0; rank <= Rank8; rank++ {
			squareIndex := FileRankToSquareIndex(file, rank)
			gotFile, gotRank := SquareIndexToFileRank(squareIndex)
			if gotFile != file || gotRank != rank {
				t.Errorf("round-trip failed for file=%d rank=%d; got file=%d rank=%d", file, rank, gotFile, gotRank)
			}
		}
	}
}

func TestFileRankNotationRoundTrip(t *testing.T) {
	for file := 0; file <= FileH; file++ {
		for rank := 0; rank <= Rank8; rank++ {
			squareNotation := FileRankToNotation(file, rank)
			gotFile, gotRank := SquareNotationToFileRank(squareNotation)
			if gotFile != file || gotRank != rank {
				t.Errorf("round-trip failed for file=%d rank=%d; got file=%d rank=%d", file, rank, gotFile, gotRank)
			}
		}
	}
}

func TestIsSquareAttacked(t *testing.T) {
	tests := []struct {
		name             string
		inputFEN         string
		inputSquare      string
		inputAttacker    Piece
		expectedAttacked bool
	}{
		{
			"attacked by an adjacent king",
			"8/8/8/4K3/8/8/8/8 w - - 0 1",
			"e4",
			White,
			true,
		},
		{
			"attacked by a knight",
			"8/8/5N2/8/8/8/8/8 w - - 0 1",
			"e4",
			White,
			true,
		},
		{
			"attacked by a bishop sliding several squares",
			"8/8/8/8/8/8/8/B7 w - - 0 1",
			"e5",
			White,
			true,
		},
		{
			"bishop's attack is blocked by a piece in between",
			"8/8/8/8/8/2N5/8/B7 w - - 0 1",
			"e5",
			White,
			false,
		},
		{
			"attacked by a rook sliding several squares",
			"8/8/8/8/8/8/8/R7 w - - 0 1",
			"a5",
			White,
			true,
		},
		{
			"rook's attack is blocked by a piece in between",
			"8/8/8/8/8/P7/8/R7 w - - 0 1",
			"a5",
			White,
			false,
		},
		{
			"attacked by a queen diagonally",
			"8/8/8/8/8/8/8/Q7 w - - 0 1",
			"e5",
			White,
			true,
		},
		{
			"a piece of the wrong color does not attack",
			"8/8/8/8/8/8/8/B7 w - - 0 1",
			"e5",
			Black,
			false,
		},
		{
			"attacked by a white pawn",
			"8/8/8/8/8/3P4/8/8 w - - 0 1",
			"e4",
			White,
			true,
		},
		{
			"attacked by a black pawn",
			"8/8/8/3p4/8/8/8/8 w - - 0 1",
			"e4",
			Black,
			true,
		},
		{
			"empty board, square is not attacked",
			"8/8/8/8/8/8/8/8 w - - 0 1",
			"e4",
			White,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := ParseFEN(tt.inputFEN)
			file, rank := SquareNotationToFileRank(tt.inputSquare)
			square := FileRankToSquareIndex(file, rank)

			attacked := board.IsSquareAttacked(square, tt.inputAttacker)
			if attacked != tt.expectedAttacked {
				t.Errorf("IsSquareAttacked(%s, %d) = %v; want %v", tt.inputSquare, tt.inputAttacker, attacked, tt.expectedAttacked)
			}
		})
	}
}
