package board

import "math/rand"

// AllInit initialize everything
func AllInit() {
	InitSq120To64()
	InitHashKeys()
	InitFilesRanksBoard()
}

// InitFilesRanksBoard initialize arrays that hold information about which rank & file a square is on the board
func InitFilesRanksBoard() {
	// Set all square to OffBoard
	for index := 0; index < BoardSquareNum; index++ {
		FilesBoard[index] = OffBoard
		RanksBoard[index] = OffBoard
	}

	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			FilesBoard[sq] = file
			RanksBoard[sq] = rank
		}
	}
}

// InitSq120To64 Initialize board covertion arrays
func InitSq120To64() {
	// Set invalid values for all squares in 120Sq array
	for index := 0; index < BoardSquareNum; index++ {
		Sq120ToSq64[index] = 65
	}

	// Set invalid values for all squares in 64Sq array
	for index := 0; index < 64; index++ {
		Sq64ToSq120[index] = 120
	}
	// The above setup is later used for fail safe check that everything is set correctly

	sq64 := 0
	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			Sq64ToSq120[sq64] = sq
			Sq120ToSq64[sq] = sq64
			sq64++
		}
	}
}

// InitHashKeys initializes hashkeys for all pieces and possible positions, for castling rights, for side to move
func InitHashKeys() {
	for i := 0; i < 13; i++ {
		for j := 0; j < 120; j++ {
			PieceKeys[i][j] = rand.Uint64() // returns a random 64 bit number
		}
	}

	SideKey = rand.Uint64()
	for i := 0; i < 16; i++ {
		CastleKeys[i] = rand.Uint64()
	}
}
