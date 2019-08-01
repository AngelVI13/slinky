package board

import (
	"fmt"
	"time"
)

const (
	//Infinite infinite value
	Infinite int = 30000
	// IsMate mate value
	IsMate int = Infinite - MaxDepth
)

// CheckUp check if time up or interrupt from GUI
func CheckUp(info *SearchInfo) {
	// check if time up or interrupt from GUI
	// fmt.Println(elapsed, info.StopTime, elapsed.After(info.StopTime))
	elapsedTime := time.Since(info.StartTime).Seconds() * 1000 // get elapsed time in ms
	if info.TimeSet == true && elapsedTime > float64(info.StopTime) {
		info.stopped = true
	}
	// if we received something from the gui -> set stopped/quit to true
	// ReadInput(info)
}

// SearchPosition searches a given position
func SearchPosition(pos *ChessBoard, info *SearchInfo) int {
	// ... iterative deepening, search init
	// for depth = 1 to maxDepth
	// 		search with alphaBeta if you have enough time left
	// you do not search to maxDepth from the start but instead search
	// with depth 1, then 2, then 3 etc, because you first identify
	// the principle variation or the potentially good moves and in this
	// way when you search again with more depth you can easily eliminate
	// a lot of bad nodes automatically

	// if we can perform a book move, do that first, otherwise perform search
	// bestMove := GetBookMove(pos)
	// if bestMove != 0 {
	// 	PerformMove(pos, info, bestMove)
	// 	return 0
	// }

	// // do normal move search
	// bestMove = NoMove
	// bestScore := -Infinite


	// moveTime := int64(time.Since(info.StartTime).Seconds() * 1000) // the UCI protocol expects milliseconds
	// if info.GameMode == UciMode {
	// 	fmt.Printf("info score cp %d depth %d nodes %d time %d ", bestScore, currentDepth, info.nodes, moveTime)
	// } else if info.GameMode == XBoardMode && info.PostThinking == true {
	// 	moveTime *= 10
	// 	fmt.Printf("%d %d %d %d", currentDepth, bestScore, moveTime, info.nodes)
	// } else if info.PostThinking == true {
	// 	fmt.Printf("score:%d depth:%d nodes:%d time:%d(ms)", bestScore, currentDepth, info.nodes, moveTime)
	// }
	// if info.GameMode == UciMode || info.PostThinking == true {
	// 	// Print the principle variation
	// 	// todo add ability to print out PV line
	// 	// pvMoves = GetPvLine(pos, currentDepth)
	// 	// fmt.Printf("pv")
	// 	// for pvNum := 0; pvNum < pvMoves; pvNum++ {
	// 	// 	fmt.Printf(" %s", PrintMove(pos.PvArray[pvNum]))
	// 	// }
	// 	// fmt.Println()
	// 	// fmt.Printf("Ordering: %.2f\n", info.failHighFirst/info.failHigh)
	// }

	// PerformMove(pos, info, bestMove)

	return 0
}

// PerformMove performs the best found move from search or book
func PerformMove(pos *ChessBoard, info *SearchInfo, bestMove int) {
	if info.GameMode == UciMode {
		fmt.Printf("bestmove %s\n", PrintMove(bestMove))
	} else if info.GameMode == XBoardMode {
		fmt.Printf("move %s\n", PrintMove(bestMove))
		pos.MakeMove(bestMove)
	} else {
		fmt.Printf("\n\n***!! Hugo makes move %s !!***\n\n", PrintMove(bestMove))
		pos.MakeMove(bestMove)
		fmt.Println(pos)
	}
}

