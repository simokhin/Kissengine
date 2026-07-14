package main

import (
	"strings"
)

const StartFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func ParseFEN(fen string) BoardState {
	board := BoardState{}

	fenFields := strings.Fields(fen)

	position := strings.Split(fenFields[0], "/")

	for i := 0; i < len(position); i++ {
		currentRank := Rank8 - i
		var currentFile int

		for j := 0; j < len(position[i]); j++ {
			if position[i][j] >= '1' && position[i][j] <= '8' {
				currentFile += int(position[i][j] - '0')
			} else {
				switch string(position[i][j]) {
				case "r":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Rook | Black
				case "n":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Knight | Black
				case "b":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Bishop | Black
				case "q":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Queen | Black
				case "k":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = King | Black
				case "p":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Pawn | Black
				case "R":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Rook | White
				case "N":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Knight | White
				case "B":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Bishop | White
				case "Q":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Queen | White
				case "K":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = King | White
				case "P":
					board.squares[FileRankToSquareIndex(currentFile, currentRank)] = Pawn | White
				}
				currentFile++
			}
		}
	}

	if fenFields[1] == "w" {
		board.sideToMove = WhiteToMove
	} else if fenFields[1] == "b" {
		board.sideToMove = BlackToMove
	}

	return board
}
