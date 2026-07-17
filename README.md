# Kissengine

A UCI chess engine written from scratch in Go on a 0x88 board, built as a project to
actually understand how move generation, search, and evaluation work under the hood.

## Features

- 0x88 board, copy-make moves, perft-verified move generation
- Negamax + alpha-beta, iterative deepening with time management for `movetime`/`wtime`
- Transposition table (Zobrist hashing, depth-preferred replacement, configurable via
  `setoption name Hash`), verified null-move pruning, LMR, quiescence search with SEE
- Move ordering via TT move, MVV-LVA, killer moves, history heuristic
- Mate-aware scoring, repetition detection, principal variation reporting
- Evaluation: material + PST, pawn structure (doubled/isolated/passed), rook open files,
  king safety, mobility
- Async `stop`, `ucinewgame`

Also included: a small GUI to play against it locally, a `referee` tool for engine-vs-engine
matches, and a `bench` tool for perf measurement.

## Project layout

```
engine/    core engine: board, move generation, search, evaluation
uci/       UCI protocol front-end (binary)
gui/       standalone GUI to play against the engine (binary)
referee/   engine-vs-engine match runner (binary)
bench/     node/time/depth benchmarking tool (binary)
notation/  shared UCI move-notation helpers
```

## Building

```
go build -o kissengine ./uci
```

The binary speaks standard UCI and works with any UCI-compatible GUI (CuteChess, Arena,
lichess-bot, etc.). To play locally with the bundled GUI instead: `go run ./gui`.

## License

MIT — see [LICENSE](LICENSE).
