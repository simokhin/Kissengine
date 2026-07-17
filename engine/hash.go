package engine

import (
	"math/rand"
	"unsafe"
)

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

const DefaultHashSizeMB = 64

var transpositionalTable []TableEntry
var tableMask ZobristHash

// SetHashSizeMB (re)allocates the transposition table for roughly sizeMB megabytes.
// The actual entry count is rounded down to the nearest power of two so the index
// mask (size-1) trick works instead of a slower modulo.
func SetHashSizeMB(sizeMB int) {
	entrySize := int(unsafe.Sizeof(TableEntry{}))
	numEntries := sizeMB * 1024 * 1024 / entrySize

	size := 1
	for size*2 <= numEntries {
		size *= 2
	}

	transpositionalTable = make([]TableEntry, size)
	tableMask = ZobristHash(size - 1)
}

var piecesOnBoardKeys [15][128]uint64
var sideToMoveKey uint64
var castleRightsKeys [4]uint64
var castleRightsFlags = [4]CastleRights{WhiteKingSide, WhiteQueenSide, BlackKingSide, BlackQueenSide}
var enPassantFileKeys [8]uint64

func Store(entry TableEntry) {
	index := entry.zobristHash & tableMask
	existing := transpositionalTable[index]

	// Depth-preferred replacement: keep whatever is already there unless it's an unused
	// slot or the new entry was searched at least as deep — otherwise a cheap, shallow
	// entry could evict an expensive, deep one under heavy hash pressure.
	if existing.zobristHash == 0 || entry.depth >= existing.depth {
		transpositionalTable[index] = entry
	}
}

func Probe(hash ZobristHash) (TableEntry, bool) {
	var tableEntry TableEntry

	index := hash & tableMask
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
	SetHashSizeMB(DefaultHashSizeMB)

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
