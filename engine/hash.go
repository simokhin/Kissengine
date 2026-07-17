package engine

import "math/rand"

type ZobristHash uint64

type TableEntry struct {
	zobristHash ZobristHash
	depth       int
	evaluation  Evaluation
	flag        PositionTypeFlag
	bestMove    Move
}

type PositionTypeFlag int

const (
	Exact PositionTypeFlag = iota
	LowerBound
	UpperBound
)

var transpositionalTable [1 << 20]TableEntry

var piecesOnBoardKeys [15][128]uint64
var sideToMoveKey uint64
var castleRightsKeys [4]uint64
var castleRightsFlags = [4]CastleRights{WhiteKingSide, WhiteQueenSide, BlackKingSide, BlackQueenSide}
var enPassantFileKeys [8]uint64

func Store(entry TableEntry) {
	index := entry.zobristHash & (1<<20 - 1)
	transpositionalTable[index] = entry
}

func Probe(hash ZobristHash) (TableEntry, bool) {
	var tableEntry TableEntry

	index := hash & (1<<20 - 1)
	tableEntry = transpositionalTable[index]

	if tableEntry.zobristHash == hash {
		return tableEntry, true
	}

	return TableEntry{}, false
}

func ComputeHash(board BoardState) ZobristHash {
	var zobristHash ZobristHash

	if board.sideToMove == BlackToMove {
		zobristHash ^= ZobristHash(sideToMoveKey)
	}

	for i, flag := range castleRightsFlags {
		if board.castleRights&flag != 0 {
			zobristHash ^= ZobristHash(castleRightsKeys[i])
		}
	}

	for i := range board.squares {
		piece := board.squares[i]
		if piece == Empty {
			continue
		}
		zobristHash ^= ZobristHash(piecesOnBoardKeys[piece][i])
	}

	if board.enPassantSquare != NoSquare {
		file, _ := SquareIndexToFileRank(board.enPassantSquare)
		zobristHash ^= ZobristHash(enPassantFileKeys[file])
	}

	return zobristHash
}

func init() {
	sideToMoveKey = rand.Uint64()

	for piece := range piecesOnBoardKeys {
		for square := range piecesOnBoardKeys[piece] {
			piecesOnBoardKeys[piece][square] = rand.Uint64()
		}
	}

	for i := range castleRightsKeys {
		castleRightsKeys[i] = rand.Uint64()
	}

	for i := range enPassantFileKeys {
		enPassantFileKeys[i] = rand.Uint64()
	}
}
