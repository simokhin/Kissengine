package main

import (
	"MyChessEngine/engine"
	"MyChessEngine/notation"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var history []engine.ZobristHash

func printInfo(result engine.SearchResult, pv []engine.Move) {
	var scoreStr string
	switch {
	case result.Score > engine.MateThreshold:
		pliesToMate := engine.Mate - result.Score
		movesToMate := (pliesToMate + 1) / 2
		scoreStr = fmt.Sprintf("mate %d", movesToMate)
	case result.Score < -engine.MateThreshold:
		pliesToMate := engine.Mate + result.Score
		movesToMate := (pliesToMate + 1) / 2
		scoreStr = fmt.Sprintf("mate -%d", movesToMate)
	default:
		scoreStr = fmt.Sprintf("cp %d", result.Score)
	}

	pvStrs := make([]string, len(pv))
	for i, move := range pv {
		pvStrs[i] = notation.MoveToUCI(move)
	}

	fmt.Printf("info depth %d nodes %d score %s pv %s\n", result.Depth, result.Nodes, scoreStr, strings.Join(pvStrs, " "))
}

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
			fmt.Printf("option name Hash type spin default %d min 1 max 1024\n", engine.DefaultHashSizeMB)
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "setoption":
			if len(fields) < 5 {
				continue
			}
			if fields[1] == "name" && fields[2] == "Hash" && fields[3] == "value" {
				sizeMB, err := strconv.Atoi(fields[4])
				if err == nil {
					engine.SetHashSizeMB(sizeMB)
				}
			}
		case "position":
			history = nil
			if len(fields) < 2 {
				continue
			}
			switch fields[1] {
			case "startpos":
				board = engine.ParseFEN(engine.StartFen)
				history = append(history, engine.ComputeHash(board))
			case "fen":
				fen := strings.Join(fields[2:8], " ")
				board = engine.ParseFEN(fen)
				history = append(history, engine.ComputeHash(board))
			}

			for i, field := range fields {
				if field == "moves" {
					for _, moveStr := range fields[i+1:] {
						move, ok := notation.ParseMove(moveStr, board)
						if ok {
							board = engine.MakeMove(board, move)
							history = append(history, engine.ComputeHash(board))
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
				result := engine.FindBestMove(board, depth, history)
				pv := engine.ExtractPV(board, result.Depth)
				board = engine.MakeMove(board, result.Move)
				history = append(history, engine.ComputeHash(board))
				printInfo(result, pv)
				fmt.Println("bestmove " + notation.MoveToUCI(result.Move))
			case "movetime":
				if len(fields) < 3 {
					continue
				}

				if len(engine.GenerateLegalMoves(board)) == 0 {
					fmt.Println("bestmove 0000")
					continue
				}

				ms, _ := strconv.Atoi(fields[2])
				result := engine.FindBestMoveByTime(board, time.Duration(ms)*time.Millisecond, history, false)
				pv := engine.ExtractPV(board, result.Depth)
				board = engine.MakeMove(board, result.Move)
				history = append(history, engine.ComputeHash(board))
				printInfo(result, pv)
				fmt.Println("bestmove " + notation.MoveToUCI(result.Move))
			case "wtime":
				if len(engine.GenerateLegalMoves(board)) == 0 {
					fmt.Println("bestmove 0000")
					continue
				}

				var wtime, btime, winc, binc int
				for i := 1; i+1 < len(fields); i += 2 {
					value, _ := strconv.Atoi(fields[i+1])
					switch fields[i] {
					case "wtime":
						wtime = value
					case "btime":
						btime = value
					case "winc":
						winc = value
					case "binc":
						binc = value
					}
				}

				myTime, myInc := wtime, winc
				if board.SideToMove().Color() == engine.Black {
					myTime, myInc = btime, binc
				}

				allocated := min(myTime/30+myInc, myTime/2)

				result := engine.FindBestMoveByTime(board, time.Duration(allocated)*time.Millisecond, history, true)
				pv := engine.ExtractPV(board, result.Depth)
				board = engine.MakeMove(board, result.Move)
				history = append(history, engine.ComputeHash(board))
				printInfo(result, pv)
				fmt.Println("bestmove " + notation.MoveToUCI(result.Move))
			}
		case "quit":
			return
		}
	}
}
