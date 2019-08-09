package board

import "fmt"

// Board general board interface that supports MCTS (UCT) algorithm
type Board interface {
	MakeMove(move int)
	TakeMove()
	GetMoves() []int               // this should be some other type
	GetResult(playerJM int) Result // this should also be some other type
	GetPlayerJustMoved() int
	GetEnemy(playerJM int) int
	String() string
}

// ChessBoard structure
type ChessBoard struct {
	Pieces     [BoardSquareNum]int
	kingSquare [2]int             // White's & black's king position
	Side       int                // which side's turn it is
	enPas      int                // square in which en passant capture is possible
	fiftyMove  int                // how many moves from the fifty move rule have been made
	histPly    int                // how many half moves have been made
	castlePerm int                // castle permissions
	posKey     uint64             // position key is a unique key stored for each position (used to keep track of 3fold repetition)
	pieceNum   [13]int            // how many pieces of each type are there currently on the board
	history    [MaxGameMoves]Undo // array that stores current position and variables before a move is made

	PlayerJustMoved int // At the root pretend the player just moved is Black i.e. White has the first move
}

func CreateBoard() ChessBoard {
	return ChessBoard{
		PlayerJustMoved: Black,
	}
}

// Reset resets a given board
func (pos *ChessBoard) Reset() {
	// Set all board positions to OffBoard
	for i := 0; i < BoardSquareNum; i++ {
		pos.Pieces[i] = OffBoard
	}

	// Set all real board positions to Empty
	for i := 0; i < InnerSquareNum; i++ {
		pos.Pieces[Sq64ToSq120[i]] = Empty
	}

	// Reset piece number
	for i := 0; i < 13; i++ {
		pos.pieceNum[i] = 0
	}

	pos.kingSquare[White] = NoSquare
	pos.kingSquare[Black] = NoSquare

	pos.Side = Both
	pos.enPas = NoSquare
	pos.fiftyMove = 0
	pos.histPly = 0
	pos.castlePerm = 0
	pos.posKey = 0
}

// GetEnemy Returns the enemy of a given player
func (pos *ChessBoard) GetEnemy(playerJM int) int {
	return playerJM ^ 1
}

// GetPlayerJustMoved returns the player that just made a move
func (pos *ChessBoard) GetPlayerJustMoved() int {
	return pos.PlayerJustMoved
}

// UpdateListsMaterial updates all material related piece lists
func (pos *ChessBoard) UpdateListsMaterial() {
	for index := 0; index < BoardSquareNum; index++ {
		piece := pos.Pieces[index]
		if piece != OffBoard && piece != Empty {
			colour := PieceColour[piece]

			pos.pieceNum[piece]++ // increment piece number

			if piece == WhiteKing || piece == BlackKing {
				pos.kingSquare[colour] = index
			}
		}
	}
}

// abs local method to compute absolute value of int without needing to convert to float
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MaterialDraw Determines if given the available pieces the position is a material draw, based on sjeng
func (pos *ChessBoard) MaterialDraw() bool {
	if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 0 &&
		pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[BlackBishop] == 0 && pos.pieceNum[WhiteBishop] == 0 {
			if pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[BlackKnight] < 3 {
				return true
			}
		} else if pos.pieceNum[WhiteKnight] == 0 && pos.pieceNum[BlackKnight] == 0 {
			if abs(pos.pieceNum[WhiteBishop]-pos.pieceNum[BlackBishop]) < 2 {
				return true
			}
		} else if (pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[WhiteBishop] == 0) ||
			(pos.pieceNum[BlackBishop] == 1 && pos.pieceNum[WhiteKnight] == 0) {
			if (pos.pieceNum[BlackKnight] < 3 && pos.pieceNum[BlackBishop] == 0) ||
				(pos.pieceNum[BlackBishop] == 1 && pos.pieceNum[BlackKnight] == 0) {
				return true
			}
		}
	} else if pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[WhiteRook] == 1 && pos.pieceNum[BlackRook] == 1 {
			if (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) < 2 &&
				(pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) < 2 {
				return true
			}
		} else if pos.pieceNum[WhiteRook] == 1 && pos.pieceNum[BlackRook] == 0 {
			if (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 0 &&
				((pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 1 ||
					(pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 2) {
				return true
			}
		} else if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 1 {
			if (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 0 &&
				((pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 1 ||
					(pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 2) {
				return true
			}
		}
	}
	return false
}

func (pos *ChessBoard) String() string {
	line := fmt.Sprintf("\nGame Board:\n\n")

	for rank := Rank8; rank >= Rank1; rank-- {
		line += fmt.Sprintf("%d  ", rank+1)
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			piece := pos.Pieces[sq]
			line += fmt.Sprintf("%3s", PieceChar[piece])
		}
		line += fmt.Sprintf("\n")
	}

	line += fmt.Sprintf("\n   ")
	for file := FileA; file <= FileH; file++ {
		line += fmt.Sprintf("%3c", 'a'+file)
	}
	line += fmt.Sprintf("\n")
	line += fmt.Sprintf("side:%c\n", SideChar[pos.Side])
	line += fmt.Sprintf("enPas:%d\n", pos.enPas)

	// Compute castling permissions
	wKCA := "-"
	if pos.castlePerm&WhiteKingCastling != 0 {
		wKCA = "K"
	}

	wQCA := "-"
	if pos.castlePerm&WhiteQueenCastling != 0 {
		wQCA = "Q"
	}

	bKCA := "-"
	if pos.castlePerm&BlackKingCastling != 0 {
		bKCA = "k"
	}

	bQCA := "-"
	if pos.castlePerm&BlackQueenCastling != 0 {
		bQCA = "q"
	}

	line += fmt.Sprintf("castle:%s%s%s%s\n", wKCA, wQCA, bKCA, bQCA)
	line += fmt.Sprintf("PosKey:%X\n", pos.posKey)
	return line
}

// GetThreeFoldRepetitionCount Detects how many repetitions for a given position
func (pos *ChessBoard) GetThreeFoldRepetitionCount() int {
	r := 0

	for i := 0; i < pos.histPly; i++ {
		if pos.history[i].posKey == pos.posKey {
			r++
		}
	}
	return r
}

// IsPositionDraw determine if position is a draw
func (pos *ChessBoard) IsPositionDraw() bool {
	// if there are pawns on the board the one of the sides can get mated
	if pos.pieceNum[WhitePawn] != 0 || pos.pieceNum[BlackPawn] != 0 {
		return false
	}
	// if there are major pieces on the board the one of the sides can get mated
	if pos.pieceNum[WhiteQueen] != 0 || pos.pieceNum[BlackQueen] != 0 || pos.pieceNum[WhiteRook] != 0 || pos.pieceNum[BlackRook] != 0 {
		return false
	}
	if pos.pieceNum[WhiteBishop] > 1 || pos.pieceNum[BlackBishop] > 1 {
		return false
	}
	if pos.pieceNum[WhiteKnight] > 1 || pos.pieceNum[BlackKnight] > 1 {
		return false
	}
	if pos.pieceNum[WhiteKnight] != 0 && pos.pieceNum[WhiteBishop] != 0 {
		return false
	}
	if pos.pieceNum[BlackKnight] != 0 && pos.pieceNum[BlackBishop] != 0 {
		return false
	}

	return true
}

// GetResult is called everytime a move is made this function is called to check if the game is over
func (pos *ChessBoard) GetResult(playerJM int) Result {

	if pos.fiftyMove > 100 {
		// fmt.Printf("1/2-1/2 {fifty move rule (claimed by Slinky)}\n")
		return Draw
	}

	if pos.GetThreeFoldRepetitionCount() >= 2 {
		// fmt.Printf("1/2-1/2 {3-fold repetition (claimed by Slinky)}\n")
		return Draw
	}

	if pos.IsPositionDraw() == true {
		// fmt.Printf("1/2-1/2 {insufficient material (claimed by Slinky)}\n")
		return Draw
	}

	if len(pos.GetMoves()) != 0 {
		return NoWinner
	}

	InCheck := pos.IsSquareAttacked(pos.kingSquare[pos.Side], pos.Side^1)

	if InCheck == true {
		if pos.Side == playerJM { // if i am the side in mate -> loss, else win
			// fmt.Printf("0-1 {black mates (claimed by Slinky)}\n")
			return Loss
		}
		// fmt.Printf("0-1 {white mates (claimed by Slinky)}\n")
		return Win
	}
	// not in check but no legal moves left -> stalemate
	// fmt.Printf("\n1/2-1/2 {stalemate (claimed by Slinky)}\n")
	return Draw
}
