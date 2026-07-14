package main

// BoardState representation
type BoardState struct {
	squares    [128]int8
	sideToMove SideToMove
}

type SideToMove int

// Pieces
const (
	Empty = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// Colors
const (
	White = 0
	Black = 8
)

// Side to move
const (
	WhiteToMove SideToMove = iota
	BlackToMove
)

// Files
const (
	FileA = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

// Ranks
const (
	Rank1 = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

func SquareIndexToFileRank(index int) (file, rank int) {
	file = index & 7
	rank = index >> 4
	return file, rank
}

func FileRankToSquareIndex(file, rank int) (squareIndex int) {
	squareIndex = rank*16 + file
	return squareIndex
}

func FileRankToNotation(file, rank int) string {
	fileLetter := string(rune('a' + file))
	rankDigit := string(rune('1' + rank))

	squareNotation := fileLetter + rankDigit

	return squareNotation
}

func SquareNotationToFileRank(squareNotation string) (file, rank int) {
	file = int(squareNotation[0] - 'a')
	rank = int(squareNotation[1] - '1')

	return file, rank
}
