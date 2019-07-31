package board

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

// CreateNewBoard returns a new instance of a board with default values
// func CreateNewBoard() TicTacToe {
// 	return TicTacToe { PlayerJustMoved: PlayerO, resultLines: getResultLines() }
// }

// MirrorBoard takes in a position and modifies it to be the mirrored version of it
func MirrorBoard(pos *Board) {
	var swapPiece = [13]int{Empty, BlackPawn, BlackKnight, BlackBishop, BlackRook, BlackQueen, BlackKing, WhitePawn, WhiteKnight, WhiteBishop, WhiteRook, WhiteQueen, WhiteKing}
	var tempPiecesSlice [64]int
	var tempSide = pos.side ^ 1
	var tempCastlePerm = 0
	var tempEnPassant = NoSquare

	if pos.castlePerm&WhiteKingCastling != 0 {
		tempCastlePerm |= BlackKingCastling
	}

	if pos.castlePerm&WhiteQueenCastling != 0 {
		tempCastlePerm |= BlackQueenCastling
	}

	if pos.castlePerm&BlackKingCastling != 0 {
		tempCastlePerm |= WhiteKingCastling
	}

	if pos.castlePerm&BlackQueenCastling != 0 {
		tempCastlePerm |= WhiteQueenCastling
	}

	if pos.enPas != NoSquare {
		tempEnPassant = Sq120(Mirror64[Sq64(pos.enPas)])
	}

	for sq := 0; sq < 64; sq++ {
		tempPiecesSlice[sq] = pos.Pieces[Sq120(Mirror64[sq])]
	}

	// clear board
	ResetBoard(pos)

	// write mirrored information to all relevant arrays
	for sq := 0; sq < 64; sq++ {
		tempPiece := swapPiece[tempPiecesSlice[sq]]
		pos.Pieces[Sq120(sq)] = tempPiece
	}

	pos.side = tempSide
	pos.castlePerm = tempCastlePerm
	pos.enPas = tempEnPassant

	pos.posKey = GeneratePosKey(pos)

	UpdateListsMaterial(pos)

	// // AssertTrue(CheckBoard(pos))
}

// MaterialDraw Determines if given the available pieces the position is a material draw, based on sjeng
func MaterialDraw(pos *Board) bool {
	if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 0 && pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[BlackBishop] == 0 && pos.pieceNum[WhiteBishop] == 0 {
			if pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[BlackKnight] < 3 {
				return true
			}
		} else if pos.pieceNum[WhiteKnight] == 0 && pos.pieceNum[BlackKnight] == 0 {
			if mathutils.Abs(pos.pieceNum[WhiteBishop]-pos.pieceNum[BlackBishop]) < 2 {
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
