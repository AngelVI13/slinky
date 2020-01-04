package board

import (
	"fmt"
	"testing"
)

// TestIsSquareAttacked checks that after init and new game is started
// for a provided move c2c4 the position is initialized correctly and KingSquares
// are computed correctly
func TestIsSquareAttacked(t *testing.T) {
	// Setup
	AllInit()
	boardState := CreateBoard()
	boardState.ParseFen(StartFen)
	moveStr := "c2c4"
	move := boardState.ParseMove(moveStr)
	if move == NoMove {
		t.Errorf("Move parsing error: %s was not parsed correclty", moveStr)
	}

	for color, kingSquare := range boardState.kingSquare {
		fmt.Printf("King square for %d is %s\n", color, PrintSquare(kingSquare))
		if kingSquare < A1 && kingSquare >= NoSquare {
			t.Errorf("King square for %d is outside of [A1; H8]", color)
		}
	}
}
