package utils

import (
	"fmt"
	"slinky/board"
	"slinky/uct"
	"strconv"
	"strings"
	"time"
)

// ConsoleLoop console loop for playing through the console
func ConsoleLoop(pos *board.ChessBoard, info *board.SearchInfo) {
	fmt.Printf("Welcome to Hugo In Console Mode!\n")
	fmt.Printf("Type help for commands\n\n")

	info.GameMode = board.ConsoleMode
	info.PostThinking = true

	moveTime := 3000 // 3 seconds move time
	engineSide := board.Both
	move := board.NoMove
	playout := false // forces engine to play until game is over

	// engineSide = board.Black
	pos.ParseFen(board.StartFen)

	command := ""

	for {
		if (pos.Side == engineSide || playout == true) && pos.GetResult(pos.PlayerJustMoved) == board.NoWinner {
			info.StartTime = time.Now()

			if moveTime != 0 {
				info.TimeSet = true
				info.StopTime = moveTime
			}

			// board.SearchPosition(pos, info)
			engineMove, _, _ := uct.GetEngineMoveFast(pos, info)
			fmt.Printf("Engine move is %s\n", board.PrintMove(engineMove))
			pos.MakeMove(engineMove)
			fmt.Println(pos)

			if playout == true {
				continue
			}
		}

		command, _ = GetInput("\nSlinky > ")
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
			fen := board.RemoveStringToTheLeftOfMarker(command, startStr)
			pos.ParseFen(fen)
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
			// Depth is not supported for MCTS implementation
			continue
		}

		if strings.Contains(command, "time") {
			moveTimeStr1 := board.RemoveStringToTheLeftOfMarker(command, "time ")
			moveTimeStr2 := board.RemoveStringToTheRightOfMarker(moveTimeStr1, " ")
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
