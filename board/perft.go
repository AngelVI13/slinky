package board

import (
	"fmt"
	"time"
)

// LeafNodes global leaf nodes counter
var LeafNodes int64

func perft(depth int, pos *ChessBoard) {
	if depth == 0 {
		LeafNodes++
		return
	}

	var moveList MoveList
	pos.GenerateAllMoves(&moveList)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		if !pos.MakeMove(moveList.Moves[moveNum]) {
			continue
		}
		perft(depth-1, pos)
		pos.TakeMove()
	}

	return
}

// PerftTest run perft test
func PerftTest(depth int, pos *ChessBoard) {

	// // AssertTrue(// CheckBoard(pos))

	fmt.Println(pos)
	fmt.Printf("\nStarting Test To Depth:%d\n", depth)
	LeafNodes = 0

	start := time.Now()

	var moveList MoveList
	pos.GenerateAllMoves(&moveList)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		move := moveList.Moves[moveNum]
		if !pos.MakeMove(move) {
			continue
		}
		cumulativeNodes := LeafNodes
		perft(depth-1, pos)
		pos.TakeMove()
		oldNodes := LeafNodes - cumulativeNodes
		fmt.Printf("move %d : %s : %d\n", moveNum+1, PrintMove(move), oldNodes)
	}

	elapsed := time.Since(start)

	fmt.Printf("\nTest Complete : %d nodes visited in: %s \n", LeafNodes, elapsed)

	return
}
