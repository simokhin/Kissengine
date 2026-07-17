package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: bench <engine> <maxDepth>")
		os.Exit(1)
	}

	enginePath := os.Args[1]
	maxDepth, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("maxDepth must be a number")
		os.Exit(1)
	}

	cmd := exec.Command(enginePath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(stdout)

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintln(stdin, "uci")
	for scanner.Scan() {
		if scanner.Text() == "uciok" {
			break
		}
	}

	fmt.Fprintln(stdin, "isready")
	for scanner.Scan() {
		if scanner.Text() == "readyok" {
			break
		}
	}

	fmt.Printf("%-6s %-12s %-10s %-12s\n", "depth", "nodes", "time", "nps")

	for depth := 1; depth <= maxDepth; depth++ {
		fmt.Fprintln(stdin, "position startpos")

		start := time.Now()
		fmt.Fprintf(stdin, "go depth %d\n", depth)

		var nodes int
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "info") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "nodes" && i+1 < len(fields) {
						nodes, _ = strconv.Atoi(fields[i+1])
					}
				}
			}
			if strings.HasPrefix(line, "bestmove") {
				break
			}
		}
		elapsed := time.Since(start)

		nps := float64(nodes) / elapsed.Seconds()
		fmt.Printf("%-6d %-12d %-10s %-12.0f\n", depth, nodes, elapsed.Round(time.Millisecond), nps)
	}

	fmt.Fprintln(stdin, "quit")
	cmd.Wait()
}
