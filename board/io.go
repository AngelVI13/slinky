package board

import (
	"fmt"
)

// PrintSquare get algebraic notation of square i.e. b2, a6 from array index
func PrintSquare(sq int) string {
	file := FilesBoard[sq]
	rank := RanksBoard[sq]

	// "a"[0] -> returns the byte value of the char 'a' -> convert to int to get ascii value
	// then add the file/rank value to it and convert back to string
	// therefore this automatically translates the files from 0-7 to a-h
	fileStr := string(int("a"[0]) + file)
	rankStr := string(int("1"[0]) + rank)

	squareStr := fileStr + rankStr
	return squareStr
}

// PrintMove prints move in algebraic notation
func PrintMove(move int) string {
	fileFrom := FilesBoard[FromSq(move)]
	rankFrom := RanksBoard[FromSq(move)]
	fileTo := FilesBoard[ToSq(move)]
	rankTo := RanksBoard[ToSq(move)]

	promoted := Promoted(move)

	moveStr := string(int("a"[0])+fileFrom) + string(int("1"[0])+rankFrom) +
		string(int("a"[0])+fileTo) + string(int("1"[0])+rankTo)

	// if this move is a promotion, add char of the piece we promote to at the end of the move string
	// i.e. if a7a8q -> we promote to Queen
	if promoted != 0 {
		pieceChar := "q"
		if IsPieceKnight[promoted] {
			pieceChar = "n"
		} else if IsPieceRookQueen[promoted] && !IsPieceBishopQueen[promoted] {
			pieceChar = "r"
		} else if !IsPieceRookQueen[promoted] && IsPieceBishopQueen[promoted] {
			pieceChar = "b"
		}
		moveStr += pieceChar
	}

	return moveStr
}

// PrintMoveList prints move list
func PrintMoveList(moveList *MoveList) {
	fmt.Println("MoveList:\n", moveList.Count)

	for index := 0; index < moveList.Count; index++ {

		move := moveList.Moves[index]

		fmt.Printf("Move:%d > %s\n", index+1, PrintMove(move))
	}
	fmt.Printf("MoveList Total %d Moves:\n\n", moveList.Count)
}
