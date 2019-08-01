package board

// --- Hashing 'macros' ---
func (pos *ChessBoard) hashPiece(piece, sq int) {
	pos.posKey ^= PieceKeys[piece][sq]
}

func (pos *ChessBoard) hashCastlePerm() {
	pos.posKey ^= CastleKeys[pos.castlePerm]
}

func (pos *ChessBoard) hashSide() {
	pos.posKey ^= SideKey
}

func (pos *ChessBoard) hashEnPass() {
	pos.posKey ^= PieceKeys[Empty][pos.enPas]
}

// ------------------------

// CastlePerm used to simplify hashing castle permissions
// Everytime we make a move we will take pos.castlePerm &= CastlePerm[sq]
// in this way if any of the rooks or the king moves, the castle permission will be
// disabled for that side. In any other move, the castle permissions will remain the
// same, since 15 is the max number associated with all possible castling permissions
// for both sides
var CastlePerm = [BoardSquareNum]int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

func (pos *ChessBoard) clearPiece(sq int) {
	pce := pos.Pieces[sq]
	pos.hashPiece(pce, sq)
	pos.Pieces[sq] = Empty
	pos.pieceNum[pce]--
}

func (pos *ChessBoard) addPiece(sq, pce int) {
	pos.hashPiece(pce, sq)
	pos.Pieces[sq] = pce
	pos.pieceNum[pce]++
}

func (pos *ChessBoard) movePiece(from, to int) {
	pce := pos.Pieces[from]

	// hash the piece out of the from square and then later hash it back in to the new square
	pos.hashPiece(pce, from)
	pos.Pieces[from] = Empty

	pos.hashPiece(pce, to)
	pos.Pieces[to] = pce
}

// MakeMove perform a move
// return false if the side to move has left themselves in check after the move i.e. illegal move
func (pos *ChessBoard) MakeMove(move int) bool {
	from := FromSq(move)
	to := ToSq(move)
	side := pos.side

	// Store has value before we do any hashing in/out of pieces etc
	pos.history[pos.histPly].posKey = pos.posKey

	// if this is an en passant move
	if move&MoveFlagEnPass != 0 {
		// if the side thats making the capture is white
		// then we need to remove the black pawn right behind the new position of the white piece
		// i.e. new_pos - 10 -> translated to array index
		if side == White {
			pos.clearPiece(to-10)
		} else {
			pos.clearPiece(to+10)
		}
	} else if move&MoveFlagCastle != 0 {
		// if its a castling move, based on the TO square, make the appopriate move, otherwise assert false
		switch to {
		case C1:
			pos.movePiece(A1, D1)
		case C8:
			pos.movePiece(A8, D8)
		case G1:
			pos.movePiece(H1, F1)
		case G8:
			pos.movePiece(H8, F8)
		default:
		}
	}

	// If the current enpassant square is SET, then we hash in the poskey
	if pos.enPas != NoSquare {
		pos.hashEnPass()
	}
	pos.hashCastlePerm() // hash out the castling permissions

	// store information to the history array about this move
	pos.history[pos.histPly].move = move
	pos.history[pos.histPly].fiftyMove = pos.fiftyMove
	pos.history[pos.histPly].enPas = pos.enPas
	pos.history[pos.histPly].castlePerm = pos.castlePerm

	// if a rook or king has moved the remove the respective castling permission from castlePerm
	pos.castlePerm &= CastlePerm[from]
	pos.castlePerm &= CastlePerm[to]
	pos.enPas = NoSquare // set enpassant square to no square

	pos.hashCastlePerm() // hash back in the castling perm

	pos.fiftyMove++ // increment firfty move rule

	// get what piece, if any, was captured in the move and if somethig was actually captured
	// i.e. captured piece is not empty remove captured piece and reset fifty move rule
	if captured := Captured(move); captured != Empty {
		pos.clearPiece(to)
		pos.fiftyMove = 0
	}

	// increase halfmove counter and ply counter values
	pos.histPly++

	// check if we need to set a new en passant square i.e. if this is a pawn start
	// then depending on the side find the piece just behind the new pawn destination
	// i.e. A4 -> compute A3 and set that as a possible enpassant capture square
	if IsPiecePawn[pos.Pieces[from]] {
		pos.fiftyMove = 0
		if move&MoveFlagPawnStart != 0 {
			if side == White {
				pos.enPas = from + 10
			} else {
				pos.enPas = from - 10
			}
			pos.hashEnPass() // hash in the enpass
		}
	}

	pos.movePiece(from, to)

	// get promoted piece and if its not empty, clear old piece (pawn)
	// and add new piece (whatever was the selected promotion piece)
	if promotedPiece := Promoted(move); promotedPiece != Empty {
		pos.clearPiece(to)
		pos.addPiece(to, promotedPiece)
	}

	// if we move the king -> update king square
	if IsPieceKing[pos.Pieces[to]] {
		pos.kingSquare[pos.side] = to
	}

	pos.side ^= 1 // change side to move
	pos.hashSide() // hash in the new side

	// check if after this move, our king is in check -> if yes -> illegal move
	if pos.IsSquareAttacked(pos.kingSquare[side], pos.side) {
		pos.TakeMove()
		return false
	}

	return true
}

// TakeMove revert move, opposite to MakeMove()
func (pos *ChessBoard) TakeMove() {
	pos.histPly--

	move := pos.history[pos.histPly].move
	from := FromSq(move)
	to := ToSq(move)

	if pos.enPas != NoSquare {
		pos.hashEnPass()
	}
	pos.hashCastlePerm()

	pos.castlePerm = pos.history[pos.histPly].castlePerm
	pos.fiftyMove = pos.history[pos.histPly].fiftyMove
	pos.enPas = pos.history[pos.histPly].enPas

	if pos.enPas != NoSquare {
		pos.hashEnPass()
	}
	pos.hashCastlePerm()

	pos.side ^= 1
	pos.hashSide()

	if MoveFlagEnPass&move != 0 {
		if pos.side == White {
			pos.addPiece(to-10, BlackPawn)
		} else {
			pos.addPiece(to+10, WhitePawn)
		}
	} else if MoveFlagCastle&move != 0 {
		switch to {
		case C1:
			pos.movePiece(D1, A1)
		case C8:
			pos.movePiece(D8, A8)
		case G1:
			pos.movePiece(F1, H1)
		case G8:
			pos.movePiece(F8, H8)
		default:
		}
	}

	pos.movePiece(to, from)

	if IsPieceKing[pos.Pieces[from]] {
		pos.kingSquare[pos.side] = from
	}

	if captured := Captured(move); captured != Empty {
		pos.addPiece(to, captured)
	}

	if promoted := Promoted(move); promoted != Empty {
		pos.clearPiece(from)
		if PieceColour[Promoted(move)] == White {
			pos.addPiece(from, WhitePawn)
		} else {
			pos.addPiece(from, BlackPawn)
		}
	}
}
