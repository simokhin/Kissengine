package main

var castleRightsMask [128]CastleRights

func init() {
	for i := range castleRightsMask {
		castleRightsMask[i] = WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide
	}

	e1 := FileRankToSquareIndex(SquareNotationToFileRank("e1"))
	a1 := FileRankToSquareIndex(SquareNotationToFileRank("a1"))
	h1 := FileRankToSquareIndex(SquareNotationToFileRank("h1"))
	e8 := FileRankToSquareIndex(SquareNotationToFileRank("e8"))
	a8 := FileRankToSquareIndex(SquareNotationToFileRank("a8"))
	h8 := FileRankToSquareIndex(SquareNotationToFileRank("h8"))

	castleRightsMask[e1] &^= WhiteKingSide | WhiteQueenSide
	castleRightsMask[a1] &^= WhiteQueenSide
	castleRightsMask[h1] &^= WhiteKingSide
	castleRightsMask[e8] &^= BlackKingSide | BlackQueenSide
	castleRightsMask[a8] &^= BlackQueenSide
	castleRightsMask[h8] &^= BlackKingSide
}
