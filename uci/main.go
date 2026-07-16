package main

import (
	"MyChessEngine/engine"
	"MyChessEngine/notation"
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
						move, ok := notation.ParseMove(moveStr, board)
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

				if len(engine.GenerateLegalMoves(board)) == 0 {
					fmt.Println("bestmove 0000")
					continue
				}

				depth, _ := strconv.Atoi(fields[2])
				bestMove := engine.FindBestMove(board, depth)
				board = engine.MakeMove(board, bestMove)
				fmt.Println("bestmove " + notation.MoveToUCI(bestMove))
			}
		case "quit":
			return
		}
	}
}

