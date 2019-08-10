package utils

import (
	"../board"
	"../uct"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func runEngine(pos *board.ChessBoard, info *board.SearchInfo, moveTime int) (gameOver bool) {
	if pos.GetResult(pos.PlayerJustMoved) == board.NoWinner {
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
		fmt.Println(fmt.Sprintf("FEN: %s\n", pos.GenerateFen()))
		return false
	}
	return true
}

// CommandLoop command loop for interacting with pyslinky GUI
func CommandLoop(pos *board.ChessBoard, info *board.SearchInfo, cmds []string) {
	info.GameMode = board.ConsoleMode
	info.PostThinking = true

	moveTime := 3000 // 3 seconds move time
	move := board.NoMove

	pos.ParseFen(board.StartFen)

	command := ""

	for _, value := range cmds {
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
			fmt.Printf("getfen - print fen of current position")
			fmt.Printf("** note ** - to reset time and depth, set to 0\n")
			fmt.Printf("enter moves using b7b8q notation\n\n\n")
			continue
		}

		if strings.Contains(command, "setboard") {
			startStr := "setboard "
			fen := board.RemoveStringToTheLeftOfMarker(command, startStr)
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
			for !runEngine(pos, info, moveTime) {
				winner := ""
				switch pos.GetResult(pos.PlayerJustMoved) {
				case board.Win:
					if pos.PlayerJustMoved == board.White {
						winner = "White wins!"
					} else {
						winner = "Black wins!"
					}
				case board.Loss:
					if pos.PlayerJustMoved == board.White {
						winner = "Black wins!"
					} else {
						winner = "White wins!"
					}
				default:
					winner = "Draw"
				}
				fmt.Printf("Game over: %s\n", winner)
			}
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

		if strings.Contains(command, "getfen") {
			fmt.Println(fmt.Sprintf("FEN: %s\n", pos.GenerateFen()))
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
			res := runEngine(pos, info, moveTime)
			if res == true {
				fmt.Println("Game is over")
			}
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
			pos.ParseFen(board.StartFen)
			continue
		}

		if strings.Contains(command, "go") {
			res := runEngine(pos, info, moveTime)
			if res == true {
				fmt.Println("Game is over")
			}
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
