package main

import (
	"MyChessEngine/engine"
	"MyChessEngine/notation"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	numGames   = 10
	moveTimeMs = 100
	moveCap    = 400
)

type Engine struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
}

type GameRecord struct {
	White  string   `json:"white"`
	Black  string   `json:"black"`
	Result string   `json:"result"`
	Moves  []string `json:"moves"`
}

type MatchRecord struct {
	Engine1     string       `json:"engine1"`
	Engine2     string       `json:"engine2"`
	Engine1Wins int          `json:"engine1_wins"`
	Engine2Wins int          `json:"engine2_wins"`
	Draws       int          `json:"draws"`
	Games       []GameRecord `json:"games"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: referee <engine1> <engine2>")
		os.Exit(1)
	}

	engine1Path := os.Args[1]
	engine2Path := os.Args[2]

	match := MatchRecord{Engine1: engine1Path, Engine2: engine2Path}

	for i := 0; i < numGames; i++ {
		white, err := startEngine(engine1Path)
		if err != nil {
			log.Fatal(err)
		}

		black, err := startEngine(engine2Path)
		if err != nil {
			log.Fatal(err)
		}

		engine1IsWhite := i%2 == 0
		if !engine1IsWhite {
			white, black = black, white
		}

		handshake(white)
		handshake(black)

		result, moves := playGame(white, black)

		fmt.Fprintln(white.stdin, "quit")
		white.cmd.Wait()
		fmt.Fprintln(black.stdin, "quit")
		black.cmd.Wait()

		whiteLabel, blackLabel := engine1Path, engine2Path
		if !engine1IsWhite {
			whiteLabel, blackLabel = engine2Path, whiteLabel
		}

		switch {
		case result == engine.WhiteWins && engine1IsWhite, result == engine.BlackWins && !engine1IsWhite:
			match.Engine1Wins++
		case result == engine.BlackWins && engine1IsWhite, result == engine.WhiteWins && !engine1IsWhite:
			match.Engine2Wins++
		default:
			match.Draws++
		}

		match.Games = append(match.Games, GameRecord{
			White:  whiteLabel,
			Black:  blackLabel,
			Result: resultString(result),
			Moves:  moves,
		})

		fmt.Printf("game %d: %s (white) vs %s (black) -> %s\n", i+1, whiteLabel, blackLabel, resultString(result))
	}

	data, err := json.MarshalIndent(match, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("match_result.json", data, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s: %d, %s: %d, draws: %d\n", engine1Path, match.Engine1Wins, engine2Path, match.Engine2Wins, match.Draws)
	fmt.Println("results written to match_result.json")
}

func playGame(white, black *Engine) (engine.Result, []string) {
	board := engine.ParseFEN(engine.StartFen)
	var moveHistory []string

	for {
		if len(engine.GenerateLegalMoves(board)) == 0 {
			if board.InCheck() {
				if board.SideToMove().Color() == engine.White {
					return engine.BlackWins, moveHistory
				}
				return engine.WhiteWins, moveHistory
			}
			return engine.Draw, moveHistory
		}

		if board.FiftyMovesRuleCount() >= 100 || len(moveHistory) >= moveCap {
			return engine.Draw, moveHistory
		}

		mover := white
		if board.SideToMove().Color() == engine.Black {
			mover = black
		}

		positionCmd := "position startpos"
		if len(moveHistory) > 0 {
			positionCmd += " moves " + strings.Join(moveHistory, " ")
		}
		fmt.Fprintln(mover.stdin, positionCmd)
		fmt.Fprintf(mover.stdin, "go movetime %d\n", moveTimeMs)

		var bestMoveStr string
		for mover.stdout.Scan() {
			line := mover.stdout.Text()
			if strings.HasPrefix(line, "bestmove") {
				bestMoveStr = strings.Fields(line)[1]
				break
			}
		}

		move, ok := notation.ParseMove(bestMoveStr, board)
		if !ok {
			return engine.Draw, moveHistory
		}

		board = engine.MakeMove(board, move)
		moveHistory = append(moveHistory, bestMoveStr)
	}
}

func resultString(r engine.Result) string {
	switch r {
	case engine.WhiteWins:
		return "white"
	case engine.BlackWins:
		return "black"
	default:
		return "draw"
	}
}

func handshake(e *Engine) {
	fmt.Fprintln(e.stdin, "uci")

	for e.stdout.Scan() {
		line := e.stdout.Text()
		if line == "uciok" {
			break
		}
	}
}

func startEngine(path string) (*Engine, error) {
	e := Engine{}

	cmd := exec.Command(path)
	e.cmd = cmd

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	e.stdin = stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	e.stdout = bufio.NewScanner(stdout)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &e, nil
}
