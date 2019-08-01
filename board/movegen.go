package board

/*
MoveGen(board, moveList)
	loops all pieces
		-> slider loop each dir, add move
			-> AddMove moveList->moves[moveList->count]=move, moveList->count++
*/

// GetMoveInt creates and returns a move int from given move information
func GetMoveInt(fromSq, toSq, capturePiece, promotionPiece, flag int) int {
	return (fromSq | (toSq << 7) | (capturePiece << 14) | (promotionPiece << 20) | flag)
}

// TODO: COMBINE AND PARAMETERIZE THE CODE BELOW for all moves GEN !!!!!!!!!!

// addQuietMove adds quiet move
func (pos *ChessBoard) addQuietMove(move int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(FromSq(move)))
	// // AssertTrue(SquareOnBoard(ToSq(move)))
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}

// addCaptureMove adds capture move
func (pos *ChessBoard) addCaptureMove(move int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(FromSq(move)))
	// // AssertTrue(SquareOnBoard(ToSq(move)))
	// // AssertTrue(PieceValid(Captured(move)))
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}

// addEnPassantMove adds quiet move
func (pos *ChessBoard) addEnPassantMove(move int, moveList *MoveList) {
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}

// addWhitePawnCaptureMove add white pawn capture move
func (pos *ChessBoard) addWhitePawnCaptureMove(from, to, cap int, moveList *MoveList) {
	// // AssertTrue(PieceValidEmpty(cap))
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank7 {
		// add all promotion with capture related moves
		pos.addCaptureMove(GetMoveInt(from, to, cap, WhiteQueen, 0), moveList)
		pos.addCaptureMove(GetMoveInt(from, to, cap, WhiteRook, 0), moveList)
		pos.addCaptureMove(GetMoveInt(from, to, cap, WhiteBishop, 0), moveList)
		pos.addCaptureMove(GetMoveInt(from, to, cap, WhiteKnight, 0), moveList)
	} else {
		// add normal capture moves without promotion
		pos.addCaptureMove(GetMoveInt(from, to, cap, Empty, 0), moveList)
	}
}

// addWhitePawnMove add white pawn normal moves
func (pos *ChessBoard) addWhitePawnMove(from, to int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank7 {
		// add normal promotion without capture
		pos.addQuietMove(GetMoveInt(from, to, Empty, WhiteQueen, 0), moveList)
		pos.addQuietMove(GetMoveInt(from, to, Empty, WhiteRook, 0), moveList)
		pos.addQuietMove(GetMoveInt(from, to, Empty, WhiteBishop, 0), moveList)
		pos.addQuietMove(GetMoveInt(from, to, Empty, WhiteKnight, 0), moveList)
	} else {
		pos.addQuietMove(GetMoveInt(from, to, Empty, Empty, 0), moveList)
	}
}

// addBlackPawnCaptureMove add black pawn capture move
func (pos *ChessBoard) addBlackPawnCaptureMove(from, to, cap int, moveList *MoveList) {
	// // AssertTrue(PieceValidEmpty(cap))
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank2 {
		// add all promotion with capture related moves
		pos.addCaptureMove(GetMoveInt(from, to, cap, BlackQueen, 0), moveList)
		pos.addCaptureMove(GetMoveInt(from, to, cap, BlackRook, 0), moveList)
		pos.addCaptureMove(GetMoveInt(from, to, cap, BlackBishop, 0), moveList)
		pos.addCaptureMove(GetMoveInt(from, to, cap, BlackKnight, 0), moveList)
	} else {
		// add normal capture moves without promotion
		pos.addCaptureMove(GetMoveInt(from, to, cap, Empty, 0), moveList)
	}
}

// addBlackPawnMove add black pawn normal moves
func (pos *ChessBoard) addBlackPawnMove(from, to int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank2 {
		// add normal promotion without capture
		pos.addQuietMove(GetMoveInt(from, to, Empty, BlackQueen, 0), moveList)
		pos.addQuietMove(GetMoveInt(from, to, Empty, BlackRook, 0), moveList)
		pos.addQuietMove(GetMoveInt(from, to, Empty, BlackBishop, 0), moveList)
		pos.addQuietMove(GetMoveInt(from, to, Empty, BlackKnight, 0), moveList)
	} else {
		pos.addQuietMove(GetMoveInt(from, to, Empty, Empty, 0), moveList)
	}
}

func (pos *ChessBoard) IsSquareOnBoard(sq int) bool {
	return FilesBoard[sq] != OffBoard
}

func (pos *ChessBoard) GenerateCastlingMoves(moveList *MoveList) {
	if pos.Side == White {
		// if the position allows white king castling
		// here we do not check if square G1 (final square after castling) is attacked
		// this will be handled at the end of the function where we will verify that all generated
		// moves are legal
		if (pos.castlePerm & WhiteKingCastling) != 0 {
			if pos.Pieces[F1] == Empty && pos.Pieces[G1] == Empty {
				if !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(F1, Black) {
					pos.addQuietMove(GetMoveInt(E1, G1, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}

		if (pos.castlePerm & WhiteQueenCastling) != 0 {
			if pos.Pieces[D1] == Empty && pos.Pieces[C1] == Empty && pos.Pieces[B1] == Empty {
				if !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(D1, Black) {
					pos.addQuietMove(GetMoveInt(E1, C1, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}
	} else {
		// castling
		if (pos.castlePerm & BlackKingCastling) != 0 {
			if pos.Pieces[F8] == Empty && pos.Pieces[G8] == Empty {
				if !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(F8, White) {
					pos.addQuietMove(GetMoveInt(E8, G8, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}

		if (pos.castlePerm & BlackQueenCastling) != 0 {
			if pos.Pieces[D8] == Empty && pos.Pieces[C8] == Empty && pos.Pieces[B8] == Empty {
				if !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(D8, White) {
					pos.addQuietMove(GetMoveInt(E8, C8, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}
	}
}

func (pos *ChessBoard) GeneratePawnMoves(sq int, moveList *MoveList) {
	var forwardOneSq, forwardTwoSq, captureLeftSq, captureRightSq int
	var enemy int
	var pawnRank int
	var pawnMoveHandler func(from, to int, ml *MoveList)
	var pawnCaptureMoveHandler func(from, to, cap int, ml *MoveList)

	if pos.Side == White {
		enemy = Black
		pawnRank = Rank2

		forwardOneSq, forwardTwoSq, captureLeftSq, captureRightSq = 10, 20, 9, 11
        pawnMoveHandler, pawnCaptureMoveHandler = pos.addWhitePawnMove, pos.addWhitePawnCaptureMove
	} else {
		enemy = White
		pawnRank = Rank7
		forwardOneSq, forwardTwoSq, captureLeftSq, captureRightSq = -10, -20, -9, -11
		pawnMoveHandler, pawnCaptureMoveHandler = pos.addBlackPawnMove, pos.addBlackPawnCaptureMove
	}

	// add simple pawn move forward if next sq is empty
	if pos.Pieces[sq + forwardOneSq] == Empty {
		pawnMoveHandler(sq, sq + forwardOneSq, moveList)
		// if we are on the second rank, generate a double pawn move if 4th rank sq is empty
		if RanksBoard[sq] == pawnRank && pos.Pieces[sq + forwardTwoSq] == Empty {
			// don't forget to set the flag for PAWN START
			pos.addQuietMove(GetMoveInt(sq, (sq + forwardTwoSq), Empty, Empty, MoveFlagPawnStart), moveList)
		}
	}

	// Capture to the left and right
	// check if the square that we are capturing on is on the board and that it has a black piece on it
	if pos.IsSquareOnBoard(sq + captureLeftSq) && PieceColour[pos.Pieces[sq + captureLeftSq]] == enemy {
		pawnCaptureMoveHandler(sq, sq + captureLeftSq, pos.Pieces[sq + captureLeftSq], moveList)
	}

	// check if the square that we are capturing on is on the board and that it has a black piece on it
	if pos.IsSquareOnBoard(sq + captureRightSq) && PieceColour[pos.Pieces[sq + captureRightSq]] == enemy {
		pawnCaptureMoveHandler(sq, sq + captureRightSq, pos.Pieces[sq + captureRightSq], moveList)
	}

	if pos.enPas != NoSquare {
		// check if the sq+9 square is equal to the enpassant square that we have stored in our pos
		if sq + captureLeftSq == pos.enPas {
			pos.addEnPassantMove(GetMoveInt(sq, sq + captureLeftSq, Empty, Empty, MoveFlagEnPass), moveList)
		}

		if sq + captureRightSq == pos.enPas {
			pos.addEnPassantMove(GetMoveInt(sq, sq + captureRightSq, Empty, Empty, MoveFlagEnPass), moveList)
		}
	}
}

func (pos *ChessBoard) GenerateSlidingMoves(sq, piece int, moveList *MoveList) {
	for i := 0; i< NumberOfDir[piece]; i++ {
		dir := PiececeDir[piece][i]
		targetSq := sq + dir

		for pos.IsSquareOnBoard(targetSq) == true {
			// BLACK ^ 1 == WHITE       WHITE ^ 1 == BLACK
			if pos.Pieces[targetSq] != Empty {
				if PieceColour[pos.Pieces[targetSq]] == pos.Side ^ 1 {
					pos.addCaptureMove(GetMoveInt(sq, targetSq, pos.Pieces[targetSq], Empty, 0), moveList)
				}

				break  // if we hit a non-empty square, we break from this direction
			}

			pos.addQuietMove(GetMoveInt(sq, targetSq, Empty, Empty, 0), moveList)
			targetSq += dir
		}
	}
}

func (pos *ChessBoard) GenerateNonSlidingMoves(sq, piece int, moveList *MoveList) {
	for i := 0; i< NumberOfDir[piece]; i++ {
		dir := PiececeDir[piece][i]
		targetSq := sq + dir

		if pos.IsSquareOnBoard(targetSq) == false {
			continue
		}

		if pos.Pieces[targetSq] != Empty {
			if PieceColour[pos.Pieces[targetSq]] == pos.Side ^ 1 {
				pos.addCaptureMove(GetMoveInt(sq, targetSq, pos.Pieces[targetSq], Empty, 0), moveList)
			}
			continue
		}
		pos.addQuietMove(GetMoveInt(sq, targetSq, Empty, Empty, 0), moveList)
	}
}

func (pos *ChessBoard) GenerateAllMoves(moveList *MoveList) {
	pos.GenerateCastlingMoves(moveList)

	for sq := 0; sq < BoardSquareNum; sq++ {
		piece := pos.Pieces[sq]

		if piece == OffBoard || PieceColour[piece] != pos.Side {
			continue
		}

		switch piece {
		case WhitePawn, BlackPawn:
			pos.GeneratePawnMoves(sq, moveList)
		case WhiteKnight, BlackKnight, WhiteKing, BlackKing:
			pos.GenerateNonSlidingMoves(sq, piece, moveList)
		case WhiteRook, BlackRook, WhiteBishop, BlackBishop, WhiteQueen, BlackQueen:
			pos.GenerateSlidingMoves(sq, piece, moveList)
		}
	}
}
