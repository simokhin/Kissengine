package engine

type Result int8

const (
	WhiteWins Result = iota
	BlackWins
	Draw
)
