package engine

import "slices"

// attackerValues converts a list of attacking piece types (from AttackersOf) into their
// material values, so they can be sorted and consumed cheapest-first during the swap.
func attackerValues(pieces []Piece) []Evaluation {
	values := make([]Evaluation, len(pieces))
	for i, p := range pieces {
		values[i] = pieceValues[p.Type()]
	}
	return values
}

// removeValue removes one occurrence of v from values (order doesn't matter — attackers
// of equal value are interchangeable for SEE purposes).
func removeValue(values []Evaluation, v Evaluation) []Evaluation {
	for i, val := range values {
		if val == v {
			return append(values[:i], values[i+1:]...)
		}
	}
	return values
}

// SEE (Static Exchange Evaluation) statically estimates the material outcome of a capture
// sequence on move.To(), assuming both sides always recapture with their cheapest available
// attacker, and stop as soon as continuing would no longer be worth it. Positive means the
// side making the initial capture comes out ahead; negative means the capture loses material.
func SEE(board BoardState, move Move) Evaluation {
	to := move.To()
	attacker := board.PieceAt(move.From())

	var attackers [2][]Evaluation // 0 = White, 1 = Black
	attackers[0] = attackerValues(board.AttackersOf(to, White))
	attackers[1] = attackerValues(board.AttackersOf(to, Black))
	slices.Sort(attackers[0])
	slices.Sort(attackers[1])

	side := 0
	if attacker.Color() == Black {
		side = 1
	}
	attackers[side] = removeValue(attackers[side], pieceValues[attacker.Type()])

	gain := []Evaluation{pieceValues[move.CapturedPiece().Type()]}
	attackerValue := pieceValues[attacker.Type()]
	side = 1 - side

	for len(attackers[side]) > 0 {
		gain = append(gain, attackerValue-gain[len(gain)-1])
		attackerValue = attackers[side][0]
		attackers[side] = attackers[side][1:]
		side = 1 - side
	}

	for i := len(gain) - 1; i > 0; i-- {
		if -gain[i] < gain[i-1] {
			gain[i-1] = -gain[i]
		}
	}

	return gain[0]
}
