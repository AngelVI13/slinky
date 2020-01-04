package board

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	// Setup
	AllInit()
	boardState := CreateBoard()
	boardState.ParseFen(StartFen)
	moveStr := "c2c4"
	move := boardState.ParseMove(moveStr)
	if move == NoMove {
		t.Errorf("Move parsing error: %s was not parsed correclty", moveStr)
	}

	boardString := boardState.String()
	if len(boardString) <= 0 {
		t.Errorf("Test")
	}
	fmt.Println(boardString)
}

func TestBoardCopy(t *testing.T) {
	boardState := CreateBoard()

	boardCopy := boardState
	boardCopy.kingSquare[0] = 7
	if boardCopy.kingSquare[0] == boardState.kingSquare[0] {
		t.Errorf("Copying failed. King squares match.")
	}

	boardCopy.history[0].fiftyMove = 33
	if boardCopy.history[0].fiftyMove == boardState.history[0].fiftyMove {
		t.Errorf("Copying failed. History items match.")
	}
}
