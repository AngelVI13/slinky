package board

import "fmt"
import "strconv"

// Board general board interface that supports MCTS (UCT) algorithm
type Board interface {
	MakeMove(move int)
	TakeMove()
	GetMoves() []int                  // this should be some other type
	GetResult(playerJM Player) Result // this should also be some other type
	GetPlayerJustMoved() Player
	GetEnemy(playerJM Player) Player
	String() string
}

// ChessBoard structure
type ChessBoard struct {
	Pieces        [BoardSquareNum]int
	kingSquare    [2]int             // White's & black's king position
	Side          int                // which side's turn it is
	enPas         int                // square in which en passant capture is possible
	fiftyMove     int                // how many moves from the fifty move rule have been made
	histPly       int                // how many half moves have been made
	castlePerm    int                // castle permissions
	posKey        uint64             // position key is a unique key stored for each position (used to keep track of 3fold repetition)
	pieceNum      [13]int            // how many pieces of each type are there currently on the board
	history       [MaxGameMoves]Undo // array that stores current position and variables before a move is made

	PlayerJustMoved int  // At the root pretend the player just moved is Black i.e. White has the first move
}

func CreateBoard() ChessBoard {
	return ChessBoard {
		PlayerJustMoved: Black,
	}
}

// CreateNewBoard returns a new instance of a board with default values
// func CreateNewBoard() TicTacToe {
// 	return TicTacToe { PlayerJustMoved: PlayerO, resultLines: getResultLines() }
// }

// Reset resets a given board
func (pos *ChessBoard) Reset() {
	// Set all board positions to OffBoard
	for i := 0; i < BoardSquareNum; i++ {
		pos.Pieces[i] = OffBoard
	}

	// Set all real board positions to Empty
	for i := 0; i < 64; i++ {
		pos.Pieces[Sq120(i)] = Empty
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

// ParseFen parse fen position string and setup a position accordingly
// TODO FIX ERROR HANDLING, now it simply returns a non-zero int whenever there is an error
func (pos *ChessBoard) ParseFen(fen string) int {
	// // AssertTrue(fen != "")
	// // AssertTrue(pos != nil)

	rank := Rank8 // we start from rank 8 since the notation starts from rank 8
	file := FileA
	piece := 0
	count := 0 // number of empty squares declared inside fen string
	sq64 := 0
	sq120 := 0

	pos.Reset()
	char := 0

	for (rank >= Rank1) && char < len(fen) {
		count = 1
		switch t := string(fen[char]); t {
		case "p", "r", "n", "b", "k", "q", "P", "R", "N", "B", "K", "Q":
			// If we have a piece related char -> set the piece to corresponding value, i.e p -> BlackPawn
			piece = PieceNotationMap[t]
		case "1", "2", "3", "4", "5", "6", "7", "8":
			// otherwise it must be a count of a number of empty squares
			piece = Empty
			count, _ = strconv.Atoi(t) // get number of empty squares and store in count
		case "/", " ":
			// if we have / or space then we are either at the end of the rank or at the end of the piece list
			// -> reset variables and continue the while loop
			rank--
			file = FileA
			char++
			continue
		default:
			fmt.Println("FEN error")
			return -1
		}

		// This loop, skips over all empty positions in a rank
		// When it comes to a piece that is different that "1"-"8" it places it on the corresponding square
		for i := 0; i < count; i++ {
			sq64 = rank*8 + file
			sq120 = Sq120(sq64)
			if piece != Empty {
				pos.Pieces[sq120] = piece
			}
			file++
		}
		char++
	}

	newChar := ""
	// newChar should be set to the side to move part of the FEN string here
	newChar = string(fen[char])
	// // AssertTrue(newChar == "w" || newChar == "b")

	if newChar == "w" {
		pos.Side = White
		pos.PlayerJustMoved = Black
	} else {
		pos.Side = Black
		pos.PlayerJustMoved = White
	}

	// move character pointer 2 characters further and it should now point to the start of the castling permissions part of FEN
	char += 2

	// Iterate over the next 4 chars - they show if white is allowed to castle king or quenside and the same for black
	for i := 0; i < 4; i++ {
		newChar = string(fen[char])
		if newChar == " " {
			// when we hit a space, it means there are no more castling permissions => break
			break
		}
		switch newChar { // Depending on the char, enable the corresponding castling permission related bit
		case "K":
			pos.castlePerm |= WhiteKingCastling
		case "Q":
			pos.castlePerm |= WhiteQueenCastling
		case "k":
			pos.castlePerm |= BlackKingCastling
		case "q":
			pos.castlePerm |= BlackQueenCastling
		default:
			break
		}
		char++
	}

	// // AssertTrue(pos.castlePerm >= 0 && pos.castlePerm <= 15)
	// move to the en passant square related part of FEN
	char++
	newChar = string(fen[char])

	if newChar != "-" {
		file := FileNotationMap[newChar]
		char++
		rank, _ := strconv.Atoi(string(fen[char])) // get rank number
		rank--                                     // decrement rank to match our indexes, i.e. Rank1 == 0

		// // AssertTrue(file >= FileA && file <= FileH)
		// // AssertTrue(rank >= Rank1 && rank <= Rank8)

		pos.enPas = FileRankToSquare(file, rank)
	}

	pos.posKey = GeneratePosKey(pos) // generate pos key for new position

	pos.UpdateListsMaterial()

	return 0
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

// Abs local method to compute absolute value of int without needing to convert to float
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MaterialDraw Determines if given the available pieces the position is a material draw, based on sjeng
func (pos *ChessBoard) MaterialDraw() bool {
	if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 0 && pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[BlackBishop] == 0 && pos.pieceNum[WhiteBishop] == 0 {
			if pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[BlackKnight] < 3 {
				return true
			}
		} else if pos.pieceNum[WhiteKnight] == 0 && pos.pieceNum[BlackKnight] == 0 {
			if Abs(pos.pieceNum[WhiteBishop]-pos.pieceNum[BlackBishop]) < 2 {
				return true
			}
		} else if (pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[WhiteBishop] == 0) || (pos.pieceNum[BlackBishop] == 1 && pos.pieceNum[WhiteKnight] == 0) {
			if (pos.pieceNum[BlackKnight] < 3 && pos.pieceNum[BlackBishop] == 0) || (pos.pieceNum[BlackBishop] == 1 && pos.pieceNum[BlackKnight] == 0) {
				return true
			}
		}
	} else if pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[WhiteRook] == 1 && pos.pieceNum[BlackRook] == 1 {
			if (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) < 2 && (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) < 2 {
				return true
			}
		} else if pos.pieceNum[WhiteRook] == 1 && pos.pieceNum[BlackRook] == 0 {
			if (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 0 && ((pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 1 || (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 2) {
				return true
			}
		} else if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 1 {
			if (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 0 && ((pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 1 || (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 2) {
				return true
			}
		}
	}
	return false
}

// todo might need to accept value receiver
func (pos *ChessBoard) String() string {
	line := fmt.Sprintf("\nGame Board:\n\n")

	for rank := Rank8; rank >= Rank1; rank-- {
		line += fmt.Sprintf("%d  ", rank+1)
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			piece := pos.Pieces[sq]
			line += fmt.Sprintf("%3c", PieceChar[piece])
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

// ParseMove parses user move and returns the MOVE int value from the GeneratedMoves for the
// position, that matches the moveStr input. For example if moveStr = 'a2a3'
// loops over all possible moves for the position, finds that move int i.e. 1451231 and returns it
func (pos *ChessBoard) ParseMove(moveStr string) (move int) {
	// THIS COULD BE DOING BYTE COMPARISON INSTEAD OF INT COMPARISON !!!!!
	// check if files for 'from' and 'to' squares are valid i.e. between 1-8
	if moveStr[1] > "8"[0] || moveStr[1] < "1"[0] {
		return NoMove
	}

	if moveStr[3] > "8"[0] || moveStr[3] < "1"[0] {
		return NoMove
	}

	// check if ranks for 'from' and 'to' squares are valid i.e. between a-h
	if moveStr[0] > "h"[0] || moveStr[0] < "a"[0] {
		return NoMove
	}

	if moveStr[2] > "h"[0] || moveStr[2] < "a"[0] {
		return NoMove
	}

	from := FileRankToSquare(int(moveStr[0]-"a"[0]), int(moveStr[1]-"1"[0]))
	to := FileRankToSquare(int(moveStr[2]-"a"[0]), int(moveStr[3]-"1"[0]))

	// fmt.Printf("Move string: %s, from: %d to: %d\n", moveStr, from, to)

	var moveList MoveList
	pos.GenerateAllMoves(&moveList)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		move := moveList.Moves[moveNum]
		if FromSq(move) == from && ToSq(move) == to {
			promPiece := Promoted(move)
			if promPiece != Empty {
				if IsPieceRookQueen[promPiece] && !IsPieceBishopQueen[promPiece] && moveStr[4] == "r"[0] {
					return move
				} else if !IsPieceRookQueen[promPiece] && IsPieceBishopQueen[promPiece] && moveStr[4] == "b"[0] {
					return move
				} else if IsPieceRookQueen[promPiece] && IsPieceBishopQueen[promPiece] && moveStr[4] == "q"[0] {
					return move
				} else if IsPieceKnight[promPiece] && moveStr[4] == "n"[0] {
					return move
				}
				continue
			}
			// must not be a promotion -> return move
			return move
		}
	}

	return NoMove
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
		fmt.Printf("1/2-1/2 {fifty move rule (claimed by Hugo)}\n")
		return Draw
	}

	if pos.GetThreeFoldRepetitionCount() >= 2 {
		fmt.Printf("1/2-1/2 {3-fold repetition (claimed by Hugo)}\n")
		return Draw
	}

	if pos.IsPositionDraw() == true {
		fmt.Printf("1/2-1/2 {insufficient material (claimed by Hugo)}\n")
		return Draw
	}

	if len(pos.GetMoves()) != 0 {
		return NoWinner
	}

	InCheck := pos.IsSquareAttacked(pos.kingSquare[pos.Side], pos.Side^1)

	if InCheck == true {
		if pos.Side == playerJM { // if i am the side in mate -> loss, else win
			// fmt.Printf("0-1 {black mates (claimed by Hugo)}\n")
			return Loss
		}
		// fmt.Printf("0-1 {white mates (claimed by Hugo)}\n")
		return Win
	}
	// not in check but no legal moves left -> stalemate
	fmt.Printf("\n1/2-1/2 {stalemate (claimed by Hugo)}\n")
	return Draw

}

