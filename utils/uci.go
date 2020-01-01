package utils

import (
	"fmt"
	"slinky/board"
	"slinky/uct"
	"strconv"
	"strings"
	"time"
)

// ParseGo parse UCI go command
// sample go commange is below
//            -white time ms -black time  -b/w increment ms -movetime ms
// go depth 6 wtime 180000 btime 100000 binc 1000 winc 1000 movetime 1000 movestogo 40
func ParseGo(line string, info *board.SearchInfo, pos *board.ChessBoard) {
	movesToGo := 15
	moveTime := -1
	timeInt := -1
	inc := 0
	info.TimeSet = false

	if strings.Contains(line, "infinite") {

	}

	if strings.Contains(line, "binc") && pos.Side == board.Black {
		incStr1 := board.RemoveStringToTheLeftOfMarker(line, "binc ")
		incStr2 := board.RemoveStringToTheRightOfMarker(incStr1, " ")
		inc, _ = strconv.Atoi(incStr2)
	}

	if strings.Contains(line, "winc") && pos.Side == board.White {
		incStr1 := board.RemoveStringToTheLeftOfMarker(line, "winc ")
		incStr2 := board.RemoveStringToTheRightOfMarker(incStr1, " ")
		inc, _ = strconv.Atoi(incStr2)
	}

	if strings.Contains(line, "wtime") && pos.Side == board.White {
		timeStr1 := board.RemoveStringToTheLeftOfMarker(line, "wtime ")
		timeStr2 := board.RemoveStringToTheRightOfMarker(timeStr1, " ")
		fmt.Println(timeStr1, "|", timeStr2)
		timeInt, _ = strconv.Atoi(timeStr2)
	}

	if strings.Contains(line, "btime") && pos.Side == board.Black {
		timeStr1 := board.RemoveStringToTheLeftOfMarker(line, "btime ")
		timeStr2 := board.RemoveStringToTheRightOfMarker(timeStr1, " ")
		timeInt, _ = strconv.Atoi(timeStr2)
	}

	if strings.Contains(line, "movestogo") {
		movesToGoStr1 := board.RemoveStringToTheLeftOfMarker(line, "movestogo ")
		movesToGoStr2 := board.RemoveStringToTheRightOfMarker(movesToGoStr1, " ")
		fmt.Println(movesToGoStr1, "|", movesToGoStr2)
		movesToGo, _ = strconv.Atoi(movesToGoStr2)
	}

	if strings.Contains(line, "movetime") {
		moveTimeStr1 := board.RemoveStringToTheLeftOfMarker(line, "movetime ")
		moveTimeStr2 := board.RemoveStringToTheRightOfMarker(moveTimeStr1, " ")
		moveTime, _ = strconv.Atoi(moveTimeStr2)
	}

	if strings.Contains(line, "depth") {
		// Depth is not supported for MCTS implementation
	}

	if moveTime != -1 {
		timeInt = moveTime
		movesToGo = 1
	}

	info.StartTime = time.Now()

	if timeInt != -1 {
		info.TimeSet = true
		timeInt /= movesToGo
		// to be on the safe side we remove 50ms from this value
		timeInt -= 50
		stopTimeInSeconds := timeInt + inc // find stop time in miliseconds
		info.StopTime = stopTimeInSeconds
	}

	fmt.Printf("time:%d start:%s stop:%d timeset:%t\n", timeInt, info.StartTime, info.StopTime, info.TimeSet)

	SearchPosition(pos, info)
}

// SearchPosition searches a given position
func SearchPosition(pos *board.ChessBoard, info *board.SearchInfo) int {
	// ... iterative deepening, search init
	// for depth = 1 to maxDepth
	// 		search with alphaBeta if you have enough time left
	// you do not search to maxDepth from the start but instead search
	// with depth 1, then 2, then 3 etc, because you first identify
	// the principle variation or the potentially good moves and in this
	// way when you search again with more depth you can easily eliminate
	// a lot of bad nodes automatically

	// if we can perform a book move, do that first, otherwise perform search
	bestMove := board.GetBookMove(pos)
	if bestMove != 0 {
		PerformMove(pos, info, bestMove)
		return 0
	}

	// do normal move search
	bestScore := 0.0
	nodes := 0
	bestMove, bestScore, nodes = uct.GetEngineMoveFast(pos, info)

	// scale from percentage to centipawn loss/gain
	// here is bestScore from point of view of enemy ?
	if bestScore < 0.5 { // 0.5 indicates a draw i.e. balanced game
		bestScore = (0.5 - bestScore) * 100
	} else {
		bestScore = (bestScore - 0.5) * 100
	}

	moveTime := int64(time.Since(info.StartTime).Seconds() * 1000) // the UCI protocol expects milliseconds
	if info.GameMode == board.UciMode {
		fmt.Printf("info score cp %d nodes %d time %d ", int(bestScore), nodes, moveTime)
	} else if info.PostThinking == true {
		fmt.Printf("score:%d time:%d(ms)", int(bestScore), moveTime)
	}
	fmt.Printf("\n")
	if info.GameMode == board.UciMode || info.PostThinking == true {
		// Print the principle variation
		// todo add ability to print out PV line
		// pvMoves = GetPvLine(pos, currentDepth)
		// fmt.Printf("pv")
		// for pvNum := 0; pvNum < pvMoves; pvNum++ {
		// 	fmt.Printf(" %s", PrintMove(pos.PvArray[pvNum]))
		// }
		// fmt.Println()
		// fmt.Printf("Ordering: %.2f\n", info.failHighFirst/info.failHigh)
	}

	PerformMove(pos, info, bestMove)

	return 0
}

// PerformMove performs the best found move from search or book
func PerformMove(pos *board.ChessBoard, info *board.SearchInfo, bestMove int) {
	if info.GameMode == board.UciMode {
		fmt.Printf("bestmove %s\n", board.PrintMove(bestMove))
	} else {
		fmt.Printf("\n\n***!! Slinky makes move %s !!***\n\n", board.PrintMove(bestMove))
		pos.MakeMove(bestMove)
		fmt.Println(pos)
	}
}

// ParsePosition parse UCI position
// the expected formats are 'position fen **' or 'position startpos'
func ParsePosition(lineIn string, pos *board.ChessBoard) {
	if strings.Contains(lineIn, "startpos") {
		pos.ParseFen(board.StartFen)
	} else {
		if strings.Contains(lineIn, "fen") {
			startStr := "fen "
			fen := board.RemoveStringToTheLeftOfMarker(lineIn, startStr)
			pos.ParseFen(fen)
		} else {
			pos.ParseFen(board.StartFen)
		}
	}

	movesStr := "moves "
	movesIdx := strings.Index(lineIn, movesStr)
	if movesIdx != -1 {
		fullMovesStr := board.RemoveStringToTheLeftOfMarker(lineIn, movesStr)
		moveSlice := strings.Split(fullMovesStr, " ")
		for i := range moveSlice {
			move := pos.ParseMove(moveSlice[i])
			if move == board.NoMove {
				break
			}
			pos.MakeMove(move)
		}
	}
	//fmt.Println(pos)
}

// UciLoop main UCI loop
func UciLoop(pos *board.ChessBoard, info *board.SearchInfo) {
	fmt.Printf("id name %s\n", board.Name)
	fmt.Printf("id author AngelVI\n")
	fmt.Println("uciok")

	line := ""

	for {
		line, _ = GetInput("")
		if len(line) < 2 {
			continue
		}

		if strings.Contains(line, "isready") {
			fmt.Println("readyok")
			continue
		} else if strings.Contains(line, "position") {
			ParsePosition(line, pos)
		} else if strings.Contains(line, "ucinewgame") {
			ParsePosition("position startpos\n", pos)
		} else if strings.Contains(line, "go") {
			ParseGo(line, info, pos)
		} else if strings.Contains(line, "quit") {
			info.Quit = true
			break
		} else if strings.Contains(line, "uci") {
			fmt.Printf("id name %s\n", board.Name)
			fmt.Printf("id author AngelVI\n")
			fmt.Println("uciok")
		}

		if info.Quit {
			break
		}
	}
}
