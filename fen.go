package main

import (
	"strconv"
	"strings"
)

const StartFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func ParseFEN(fen string) BoardState {
	board := BoardState{}

	fenFields := strings.Fields(fen)

	position := strings.Split(fenFields[0], "/")

	for i := range position {
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

	switch fenFields[1] {
	case "w":
		board.sideToMove = WhiteToMove
	case "b":
		board.sideToMove = BlackToMove
	}

	for i := 0; i < len(fenFields[2]); i++ {
		switch fenFields[2][i] {
		case 'K':
			board.castleRights |= WhiteKingSide
		case 'Q':
			board.castleRights |= WhiteQueenSide
		case 'k':
			board.castleRights |= BlackKingSide
		case 'q':
			board.castleRights |= BlackQueenSide
		default:
			continue
		}
	}

	if fenFields[3] == "-" {
		board.enPassantSquare = NoSquare
	} else {
		board.enPassantSquare = FileRankToSquareIndex(SquareNotationToFileRank(fenFields[3]))
	}

	fiftyMovesRuleCount, _ := strconv.Atoi(fenFields[4])
	board.fiftyMovesRuleCount = fiftyMovesRuleCount

	moves, _ := strconv.Atoi(fenFields[5])
	board.movesCount = moves

	return board
}
