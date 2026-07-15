package main

type Piece int8

// Pieces
const (
	Empty Piece = iota
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

// Color shows what color the piece is (0 - white, 8 - black)
func (p Piece) Color() Piece {
	return p & Black
}

// Type shows what type the piece is
func (p Piece) Type() Piece {
	return p & 0x7
}
