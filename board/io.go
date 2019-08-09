package board

import (
	"fmt"
)

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
