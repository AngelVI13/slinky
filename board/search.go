package utils

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
func SearchPosition(pos *Board, info *SearchInfo) int {
	// ... iterative deepening, search init
	// for depth = 1 to maxDepth
	// 		search with alphaBeta if you have enough time left
	// you do not search to maxDepth from the start but instead search
	// with depth 1, then 2, then 3 etc, because you first identify
	// the principle variation or the potentially good moves and in this
	// way when you search again with more depth you can easily eliminate
	// a lot of bad nodes automatically

	// if we can perform a book move, do that first, otherwise perform search
	bestMove := GetBookMove(pos)
	if bestMove != 0 {
		PerformMove(pos, info, bestMove)
		return 0
	}

	// do normal move search
	bestMove = NoMove
	bestScore := -Infinite
	ClearForSearch(pos, info)

	for currentDepth := 1; currentDepth <= info.Depth; currentDepth++ {
		//                    *alpha     *beta
		bestScore = AlphaBeta(-Infinite, Infinite, currentDepth, pos, info, true)

		if info.stopped == true {
			break
		}

		pvMoves := GetPvLine(pos, currentDepth)
		bestMove = pos.PvArray[0]

		moveTime := int64(time.Since(info.StartTime).Seconds() * 1000) // the UCI protocol expects milliseconds
		if info.GameMode == UciMode {
			fmt.Printf("info score cp %d depth %d nodes %d time %d ", bestScore, currentDepth, info.nodes, moveTime)
		} else if info.GameMode == XBoardMode && info.PostThinking == true {
			moveTime *= 10
			fmt.Printf("%d %d %d %d", currentDepth, bestScore, moveTime, info.nodes)
		} else if info.PostThinking == true {
			fmt.Printf("score:%d depth:%d nodes:%d time:%d(ms)", bestScore, currentDepth, info.nodes, moveTime)
		}
		if info.GameMode == UciMode || info.PostThinking == true {
			// Print the principle variation
			pvMoves = GetPvLine(pos, currentDepth)
			fmt.Printf("pv")
			for pvNum := 0; pvNum < pvMoves; pvNum++ {
				fmt.Printf(" %s", PrintMove(pos.PvArray[pvNum]))
			}
			fmt.Println()
			// fmt.Printf("Ordering: %.2f\n", info.failHighFirst/info.failHigh)
		}
	}

	PerformMove(pos, info, bestMove)

	return 0
}

// PerformMove performs the best found move from search or book
func PerformMove(pos *Board, info *SearchInfo, bestMove int) {
	if info.GameMode == UciMode {
		fmt.Printf("bestmove %s\n", PrintMove(bestMove))
	} else if info.GameMode == XBoardMode {
		fmt.Printf("move %s\n", PrintMove(bestMove))
		MakeMove(pos, bestMove)
	} else {
		fmt.Printf("\n\n***!! Hugo makes move %s !!***\n\n", PrintMove(bestMove))
		MakeMove(pos, bestMove)
		PrintBoard(pos)
	}
}

// ClearForSearch clear all info for search
func ClearForSearch(pos *Board, info *SearchInfo) {
	for i := 0; i < 13; i++ {
		for j := 0; j < BoardSquareNum; j++ {
			pos.searchHistory[i][j] = 0
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < MaxDepth; j++ {
			pos.searchKillers[i][j] = 0
		}
	}

	pos.HashTable.overWrite = 0
	pos.HashTable.hit = 0
	pos.HashTable.cut = 0

	pos.ply = 0
	info.stopped = false
	info.nodes = 0
	info.failHigh = 0
	info.failHighFirst = 0
}

// IsRepetition Determine if position is a repetition
func IsRepetition(pos *Board) bool {
	// Loop through moves and check if the current position is equal to the posKey of any previous positions
	// since when we have a capture or a pawn move i.e. when we reset the 50 move counter, we can not go back
	// to the same position - we cannot unmove a pawn or uncapture a piece -> limit the search to the last time
	// the fifty move was reset
	for i := pos.histPly - pos.fiftyMove; i < pos.histPly-1; i++ {

		// // AssertTrue(i >= 0 && i < MaxGameMoves)
		if pos.posKey == pos.history[i].posKey {
			return true
		}
	}

	return false
}
