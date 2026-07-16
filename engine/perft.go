package engine

func Perft(board BoardState, depth int) uint64 {
	if depth == 0 {
		return 1
	}
	var nodes uint64
	for _, move := range GenerateLegalMoves(board) {
		newBoard := MakeMove(board, move)
		nodes += Perft(newBoard, depth-1)
	}
	return nodes
}
