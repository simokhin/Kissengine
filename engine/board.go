package engine

type BoardState struct {
	squares             [128]Piece
	sideToMove          SideToMove
	castleRights        CastleRights
	enPassantSquare     Square
	fiftyMovesRuleCount int
	movesCount          int
}

type SideToMove int

type Square int

type CastleRights int

const NoSquare = -1

// Castle rights
const (
	WhiteKingSide CastleRights = 1 << iota
	WhiteQueenSide
	BlackKingSide
	BlackQueenSide
)

// SideToMove
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

func (b BoardState) PieceAt(s Square) Piece {
	return b.squares[s]
}

func (b BoardState) SideToMove() SideToMove {
	return b.sideToMove
}

func (s SideToMove) Color() Piece {
	if s == WhiteToMove {
		return White
	}
	return Black
}

func SquareIndexToFileRank(index Square) (file, rank int) {
	file = int(index & 7)
	rank = int(index >> 4)
	return file, rank
}

func FileRankToSquareIndex(file, rank int) (squareIndex Square) {
	squareIndex = Square(rank*16 + file)
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

func (s Square) IsOnBoard() bool {
	if s&0x88 != 0 {
		return false
	}
	return true
}

func (b BoardState) IsSquareAttacked(square Square, attackerColor Piece) bool {
	attackerColor &= Black

	// By King
	for i := range KingOffsets {
		candidateSquare := square + KingOffsets[i]

		if !candidateSquare.IsOnBoard() {
			continue
		} else {
			if b.squares[candidateSquare] == attackerColor|King {
				return true
			}
		}
	}

	// By Knight
	for i := range KnightOffsets {
		candidateSquare := square + KnightOffsets[i]

		if !candidateSquare.IsOnBoard() {
			continue
		} else {
			if b.squares[candidateSquare] == attackerColor|Knight {
				return true
			}
		}
	}

	// By Bishop or Queen
	for i := range BishopOffsets {
		candidateSquare := square

		for {
			candidateSquare += BishopOffsets[i]

			if !candidateSquare.IsOnBoard() {
				break
			} else {
				piece := b.squares[candidateSquare]
				if piece != Empty {
					pieceColor := piece & Black
					if pieceColor != attackerColor {
						break
					} else {
						if piece.Type() == Bishop || piece.Type() == Queen {
							return true
						}
						break
					}
				}
			}
		}
	}

	// By Rook or Queen
	for i := range RookOffsets {
		candidateSquare := square

		for {
			candidateSquare += RookOffsets[i]

			if !candidateSquare.IsOnBoard() {
				break
			} else {
				piece := b.squares[candidateSquare]
				if piece != Empty {
					pieceColor := piece & Black
					if pieceColor != attackerColor {
						break
					} else {
						if piece.Type() == Rook || piece.Type() == Queen {
							return true
						}
						break
					}
				}
			}
		}
	}

	// By Pawn
	switch attackerColor {
	case White:
		for i := range WhitePawnAttackOffsets {
			candidateSquare := square + WhitePawnAttackOffsets[i]
			if !candidateSquare.IsOnBoard() {
				continue
			} else {
				piece := b.squares[candidateSquare]
				if piece != attackerColor|Pawn {
					continue
				} else {
					return true
				}
			}
		}
	case Black:
		for i := range BlackPawnAttackOffsets {
			candidateSquare := square + BlackPawnAttackOffsets[i]
			if !candidateSquare.IsOnBoard() {
				continue
			} else {
				piece := b.squares[candidateSquare]
				if piece != attackerColor|Pawn {
					continue
				} else {
					return true
				}
			}
		}
	}

	return false
}
