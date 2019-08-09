package board

import (
	"fmt"
	"strconv"
	"strings"
)

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

// ParseFen parse fen position string and setup a position accordingly
// TODO Split into smaller parts
func (pos *ChessBoard) ParseFen(fen string) {
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
			panic("FEN error")
		}

		// This loop, skips over all empty positions in a rank
		// When it comes to a piece that is different that "1"-"8" it places it on the corresponding square
		for i := 0; i < count; i++ {
			sq64 = rank*8 + file
			sq120 = Sq64ToSq120[sq64]
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

	if newChar == "w" {
		pos.Side = White
		pos.PlayerJustMoved = Black
	} else if newChar == "b" {
		pos.Side = Black
		pos.PlayerJustMoved = White
	} else {
		panic(fmt.Sprintf("Unknown side to move: %s", newChar))
	}

	// move char pointer 2 chars further and it should now point to the start of the castling permissions part of FEN
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

		if (file < FileA || file > FileH) || (rank < Rank1 || rank > Rank8) {
			panic(fmt.Sprintf("File or rank out of bounds: file(%d) rank(%d)", file, rank))
		}

		pos.enPas = FileRankToSquare(file, rank)
	}

	pos.posKey = GeneratePosKey(pos) // generate pos key for new position

	pos.UpdateListsMaterial()
}

func (pos *ChessBoard) GenerateFen() (fen string) {
	emptyCount := 0
	for idx, value := range pos.Pieces {
		if value == OffBoard {
			continue
		}

		if Sq120ToSq64[idx] > 0 && Sq120ToSq64[idx]%RowSize == 0 {
			if emptyCount != 0 {
				fen += strconv.Itoa(emptyCount)
				emptyCount = 0
			}
			fen += "/"
		}

		if value == Empty {
			emptyCount += 1
		} else {
			if emptyCount != 0 {
				fen += strconv.Itoa(emptyCount)
				emptyCount = 0
			}

			piece := PieceChar[value]
			fen += piece
		}
	}

	// ------------
	// -- This inverts row string to match expected format
	// current fen is described as A1B1C1 etc.., expected is A8B8C8
	fenRows := strings.Split(fen, "/")
	rowsLength := len(fenRows)
	fen = ""
	for idx := range fenRows {
		fen += fmt.Sprintf("%s/", fenRows[rowsLength-idx-1])
	}
	// ------------

	side := "w"
	if pos.Side == Black {
		side = "b"
	}
	fen += fmt.Sprintf(" %s", side)

	// Compute castling permissions
	// todo convert to method and use also inside String()
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
	fen += fmt.Sprintf(" %s%s%s%s", wKCA, wQCA, bKCA, bQCA)

	if pos.enPas != NoSquare {
		fen += fmt.Sprintf(" %s", PrintSquare(pos.enPas))
	} else {
		fen += " -"
	}
	// todo add additional parameters (halfmove clock and fullmove number)
	return fen
}
