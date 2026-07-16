package notation

import "MyChessEngine/engine"

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
