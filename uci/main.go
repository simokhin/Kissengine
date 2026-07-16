package main

import (
	"MyChessEngine/engine"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	board := engine.BoardState{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		command := fields[0]

		switch command {
		case "uci":
			fmt.Println("id name Kissengine")
			fmt.Println("id author Nikita Simokhin")
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "position":
			if len(fields) < 2 {
				continue
			}
			switch fields[1] {
			case "startpos":
				board = engine.ParseFEN(engine.StartFen)
			case "fen":
				fen := strings.Join(fields[2:8], " ")
				board = engine.ParseFEN(fen)
			}

			for i, field := range fields {
				if field == "moves" {
					for _, moveStr := range fields[i+1:] {
						move, ok := ParseMove(moveStr, board)
						if ok {
							board = engine.MakeMove(board, move)
						}
					}
				}
			}
		case "go":
			if len(fields) < 2 {
				continue
			}
			switch fields[1] {
			case "depth":
				if len(fields) < 3 {
					continue
				}
				depth, _ := strconv.Atoi(fields[2])
				bestMove := engine.FindBestMove(board, depth)
				board = engine.MakeMove(board, bestMove)
				fmt.Println("bestmove " + MoveToUCI(bestMove))
			}
		case "quit":
			return
		}
	}
}

func ParseMove(moveNotation string, board engine.BoardState) (engine.Move, bool) {
	from := moveNotation[0:2]
	to := moveNotation[2:4]

	var promotionLetter string
	if len(moveNotation) > 4 {
		promotionLetter = moveNotation[4:5]
	}

	fromSquare := engine.FileRankToSquareIndex(engine.SquareNotationToFileRank(from))
	toSquare := engine.FileRankToSquareIndex(engine.SquareNotationToFileRank(to))

	var promotionPiece engine.Piece
	switch promotionLetter {
	case "q":
		promotionPiece = engine.Queen
	case "b":
		promotionPiece = engine.Bishop
	case "r":
		promotionPiece = engine.Rook
	case "n":
		promotionPiece = engine.Knight
	default:
		promotionPiece = engine.Empty
	}

	for _, move := range engine.GenerateLegalMoves(board) {
		if move.From() == fromSquare && move.To() == toSquare && move.Promotion() == promotionPiece {
			return move, true
		}
	}

	return 0, false
}

func MoveToUCI(move engine.Move) string {
	fromFile, fromRank := engine.SquareIndexToFileRank(move.From())
	toFile, toRank := engine.SquareIndexToFileRank(move.To())

	uci := engine.FileRankToNotation(fromFile, fromRank) + engine.FileRankToNotation(toFile, toRank)

	switch move.Promotion() {
	case engine.Queen:
		uci += "q"
	case engine.Rook:
		uci += "r"
	case engine.Bishop:
		uci += "b"
	case engine.Knight:
		uci += "n"
	}

	return uci
}
