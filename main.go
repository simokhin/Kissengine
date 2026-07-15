package main

func main() {
	board := ParseFEN("8/8/7n/7K/8/8/8/8 w - - 0 1")

	from := FileRankToSquareIndex(SquareNotationToFileRank("h5"))

	GenerateJumpingPieceMoves(from, board, King)
}
