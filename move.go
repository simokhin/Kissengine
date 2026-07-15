package main

type Move int

type Flag int8

const (
	QuietMove Flag = iota
	DoublePawnMove
	EnPassantCapture
	Capture
	KingCastle
	QueenCastle
)

// NewMove creates a new move
func NewMove(from, to Square, promotion Piece, flag Flag, capturedPiece Piece) Move {
	return Move(from) | Move(to)<<7 | Move(promotion)<<14 | Move(flag)<<17 | Move(capturedPiece)<<20
}

// From shows from which square the move was made
func (m Move) From() Square {
	return Square(m & 0x7F)
}

// To shows which square the move was made to
func (m Move) To() Square {
	return Square((m >> 7) & 0x7F)
}

// Promotion shows what piece the pawn has been promoted to
func (m Move) Promotion() Piece {
	return Piece((m >> 14) & 0x7)
}

// Flag gives the parameters of the move
func (m Move) Flag() Flag {
	return Flag((m >> 17) & 0x7)
}

// CapturedPiece shows which piece was captured
func (m Move) CapturedPiece() Piece {
	return Piece((m >> 20) & 0xF)
}
