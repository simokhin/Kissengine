package engine

import (
	"fmt"
	"testing"
)

// Reference positions and node counts from https://www.chessprogramming.org/Perft_Results
func TestPerft(t *testing.T) {
	type perftCase struct {
		depth         int
		expectedNodes uint64
	}

	tests := []struct {
		name  string
		fen   string
		cases []perftCase
	}{
		{
			"starting position",
			StartFen,
			[]perftCase{
				{1, 20},
				{2, 400},
				{3, 8902},
				{4, 197281},
			},
		},
		{
			"kiwipete",
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			[]perftCase{
				{1, 48},
				{2, 2039},
				{3, 97862},
				{4, 4085603},
			},
		},
		{
			"position 3",
			"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			[]perftCase{
				{1, 14},
				{2, 191},
				{3, 2812},
				{4, 43238},
			},
		},
		{
			"position 4",
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			[]perftCase{
				{1, 6},
				{2, 264},
				{3, 9467},
				{4, 422333},
			},
		},
		{
			"position 5",
			"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			[]perftCase{
				{1, 44},
				{2, 1486},
				{3, 62379},
				{4, 2103487},
			},
		},
		{
			"position 6",
			"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			[]perftCase{
				{1, 46},
				{2, 2079},
				{3, 89890},
				{4, 3894594},
			},
		},
	}

	for _, tt := range tests {
		board := ParseFEN(tt.fen)

		for _, c := range tt.cases {
			t.Run(fmt.Sprintf("%s depth %d", tt.name, c.depth), func(t *testing.T) {
				got := Perft(board, c.depth)
				if got != c.expectedNodes {
					t.Errorf("want %d nodes, got %d nodes", c.expectedNodes, got)
				}
			})
		}
	}
}
