package engine

type Result int8

const (
	WhiteWins Result = iota
	BlackWins
	Draw
)

func PlayMatch() {

}

func PlayGame(board BoardState) Result {
	for {
		moves := GenerateLegalMoves(board)

		if len(moves) == 0 {
			if board.InCheck() {
				if board.SideToMove().Color() == White {
					return BlackWins
				} else {
					return WhiteWins
				}
			} else {
				return Draw
			}
		}

		if board.fiftyMovesRuleCount >= 100 {
			return Draw
		}

		if board.movesCount >= 200 {
			return Draw
		}

		result := FindBestMove(board, 4)
		board = MakeMove(board, result.Move)
	}
}
