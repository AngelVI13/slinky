package utils

import (
	"fmt"
	board "local/slinky/board"
	uct "local/slinky/uct"
	stringutils "local/string-utils"
	"strconv"
	"strings"
	"time"
)

// CommandLoop command loop for interacting with pyslinky GUI
func CommandLoop(pos *board.ChessBoard, info *board.SearchInfo, cmds []string) {
	info.GameMode = board.ConsoleMode
	info.PostThinking = true

	depth := board.MaxDepth
	moveTime := 3000 // 3 seconds move time
	engineSide := board.Both
	move := board.NoMove
	playout := false  // forces engine to play until game is over

	// engineSide = board.Black
	pos.ParseFen(board.StartFen)

	command := ""

	for _, value := range cmds {
		// todo this doesn't work because this is not an infinite for loop
		// todo i.e. need to move this code to under cmd "go" or "force" etc.
		if (pos.Side == engineSide || playout == true) && pos.GetResult(pos.PlayerJustMoved) == board.NoWinner {
			info.StartTime = time.Now()
			info.Depth = depth

			if moveTime != 0 {
				info.TimeSet = true
				info.StopTime = moveTime
			}

			// board.SearchPosition(pos, info)
			engineMove, _, _ := uct.GetEngineMoveFast(pos, 10000, info)
			fmt.Printf("Engine move is %s\n", board.PrintMove(engineMove))
			pos.MakeMove(engineMove)
			fmt.Println(pos)

			if playout == true {
				continue
			}
		}

		command = value
		if len(command) < 2 {
			continue
		}

		if strings.Contains(command, "help") {
			fmt.Printf("Commands:\n")
			fmt.Printf("quit - quit game\n")
			fmt.Printf("force - computer will not think\n")
			fmt.Printf("print - show board\n")
			fmt.Printf("post - show thinking\n")
			fmt.Printf("nopost - do not show thinking\n")
			fmt.Printf("new - start new game\n")
			fmt.Printf("setboard x - set position to fen x\n")
			fmt.Printf("go - set computer thinking\n")
			fmt.Printf("depth x - set depth to x\n")
			fmt.Printf("time x - set thinking time to x seconds (depth still applies if set)\n")
			fmt.Printf("view - show current depth and moveTime settings\n")
			fmt.Printf("showline - show current move line so far\n")
			fmt.Printf("getmoves - show all moves")
			fmt.Printf("playout - force engine to play position till end")
			fmt.Printf("** note ** - to reset time and depth, set to 0\n")
			fmt.Printf("enter moves using b7b8q notation\n\n\n")
			continue
		}

		if strings.Contains(command, "setboard") {
			engineSide = board.Both
			startStr := "setboard "
			fen := stringutils.RemoveStringToTheLeftOfMarker(command, startStr)
			pos.ParseFen(fen)
			continue
		}

		if strings.Contains(command, "position") {
			ParsePosition(command, pos)
			continue
		}

		if strings.Contains(command, "quit") {
			info.Quit = true
			break
		}

		if strings.Contains(command, "playout") {
			playout = true
			continue
		}

		if strings.Contains(command, "getmoves") {
			moves := pos.GetMoves()
			fmt.Printf("Moves found: %d -> ", len(moves))
			for i := 0; i < len(moves); i++ {
				fmt.Printf("%s, ", board.PrintMove(moves[i]))
			}
			fmt.Printf("\n")
			continue
		}

		if strings.Compare(command, "post") == 0 {
			info.PostThinking = true
			continue
		}

		if strings.Contains(command, "print") {
			fmt.Println(pos)
			continue
		}

		if strings.Compare(command, "nopost") == 0 {
			info.PostThinking = false
			continue
		}

		if strings.Contains(command, "force") {
			engineSide = board.Both
			continue
		}

		if strings.Contains(command, "view") {
			// todo remove depth option
			if depth == board.MaxDepth {
				fmt.Printf("depth not set ")
			} else {
				fmt.Printf("depth %d", depth)
			}

			if moveTime != 0 {
				fmt.Printf(" moveTime %ds\n", moveTime/1000)
			} else {
				fmt.Printf(" moveTime not set\n")
			}
			continue
		}

		if strings.Contains(command, "showline") {
			fmt.Println(board.GetBookMove(pos))
			continue
		}

		if strings.Contains(command, "depth") {
			depthStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "depth ")
			depthStr2 := stringutils.RemoveStringToTheRightOfMarker(depthStr1, " ")
			depth, _ = strconv.Atoi(depthStr2)
			if depth == 0 {
				depth = board.MaxDepth
			}
			continue
		}

		if strings.Contains(command, "time") {
			moveTimeStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "time ")
			moveTimeStr2 := stringutils.RemoveStringToTheRightOfMarker(moveTimeStr1, " ")
			moveTime, _ = strconv.Atoi(moveTimeStr2)
			moveTime *= 1000
			continue
		}

		if strings.Contains(command, "new") {
			engineSide = board.Black
			pos.ParseFen(board.StartFen)
			continue
		}

		if strings.Contains(command, "go") {
			engineSide = pos.Side
			continue
		}

		move = pos.ParseMove(command)
		if move == board.NoMove {
			fmt.Printf("Command unknown:%s\n", command)
			continue
		}
		pos.MakeMove(move)
	}
}
