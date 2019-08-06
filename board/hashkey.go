package board

// GeneratePosKey takes a position and calculates a unique hashkey for it
func GeneratePosKey(pos *ChessBoard) (hashKey uint64) {
	var finalKey uint64
	piece := Empty

	for sq := 0; sq < BoardSquareNum; sq++ {
		piece = pos.Pieces[sq]
		// Do not calculate hashkey for squares that are not on the actual board, i.e. have value of NoSquare
		// Also do not calculate hashkey for an empty square
		if (piece != NoSquare) && (piece != Empty) && (piece != OffBoard) {
			// Check if we have a valid piece
			// Add/remove (xor) the hash value for a given piece and for a given position from the final hash value
			finalKey ^= PieceKeys[piece][sq]
		}
	}

	if pos.Side == White {
		finalKey ^= SideKey
	}

	if pos.enPas != NoSquare {
		// We have already generated hash keys for all pieces + Empty
		// => the hashkeys for value empty are used for en passant hash calculations
		finalKey ^= PieceKeys[Empty][pos.enPas]
	}

	finalKey ^= CastleKeys[pos.castlePerm]

	return finalKey
}
