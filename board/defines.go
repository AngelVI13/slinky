package utils

import "time"

const (
	// Name is the name of the chess engine
	Name = "Hugo 1.0"
	// BoardSquareNum is the total number of squares in the board representation
	BoardSquareNum = 120
	// BookFile book filename
	BookFile = "utils/book.txt"
)

//

// Defines for piece values
const (
	Empty int = iota
	WhitePawn
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

// Defines for ranks
const (
	Rank1 int = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	RankNone
)

// Defines for files
const (
	FileA int = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
	FileNone
)

// Defines for colours
const (
	White int = iota
	Black
	Both
)

// Defines for MC simulation
const (
	PlayerWhite = 1
	PlayerBlack = -1
	NoPlayer    = 0
)

// Defines for board square indexes
const (
	// Rank 1
	A1 int = iota + 21 // iota = 0, A1 = 21
	B1                 // iota = 1
	C1                 // iota = 2
	D1                 // iota = 3
	E1                 // iota = 4
	F1                 // iota = 5
	G1                 // iota = 6
	H1                 // iota = 7
	// Rank 2
	A2 int = iota + 23 // iota = 8
	B2                 // iota = 9
	C2                 // iota = 10
	D2                 // iota = 11
	E2                 // iota = 12
	F2                 // iota = 13
	G2                 // iota = 14
	H2                 // iota = 15
	// Rank 3
	A3 int = iota + 25 // iota = 16
	B3                 // iota = 17
	C3                 // iota = 18
	D3                 // iota = 19
	E3                 // iota = 20
	F3                 // iota = 21
	G3                 // iota = 22
	H3                 // iota = 23
	// Rank 4
	A4 int = iota + 27 // 51
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	// Rank 5
	A5 int = iota + 29 // 61
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	// Rank 6
	A6 int = iota + 31 // 71
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	// Rank 7
	A7 int = iota + 33 // 81
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	// Rank 8
	A8 int = iota + 35 // 91
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	// No square
	NoSquare // 99
	OffBoard // 100
)

// Defines for castling rights
// The values are such that they each represent a bit from a 4 bit int value
// for example if white can castle kingside and black can castle queenside
// the 4 bit int value is going to be 1001
const (
	WhiteKingCastling  int = 1
	WhiteQueenCastling int = 2
	BlackKingCastling  int = 4
	BlackQueenCastling int = 8
)

const (
	// MaxGameMoves maximum number halfmoves allowed
	MaxGameMoves int = 2048
)

// Undo struct
type Undo struct {
	move       int
	castlePerm int
	enPas      int
	fiftyMove  int
	posKey     uint64
}

// Board structure
type Board struct {
	Pieces        [BoardSquareNum]int
	Pawns         [3]uint64          // number of white pawns, number of black pawns, number of both pawns
	kingSquare    [2]int             // White's & black's king position
	side          int                // which side's turn it is
	enPas         int                // square in which en passant capture is possible
	fiftyMove     int                // how many moves from the fifty move rule have been made
	ply           int                // depth of search algorithm
	histPly       int                // how many half moves have been made
	castlePerm    int                // castle permissions
	posKey        uint64             // position key is a unique key stored for each position (used to keep track of 3fold repetition)
	pieceNum      [13]int            // how many pieces of each type are there currently on the board
	bigPieceNum   [2]int             // number of big pieces on the board (anything thats not a pawn) for each colour and for both
	majorPieceNum [2]int             // number of major pieces on the board (rooks and queens) for each colour and for both
	minorPieceNum [2]int             // number of minor pieces on the board (bishops and knights) for each colour and for both
	material      [2]int             // material scores for black and white
	history       [MaxGameMoves]Undo // array that stores current position and variables before a move is made
	// pieceList contains the squares of all pieces on the board, this makes it faster to iterate and generate moves for (instead of iterating over pieces slice (too big))
	// 13 is the total number of pieces for white and black combined, 10 is the maximum possible number of each piece to occur in a game
	pieceList [13][10]int
	HashTable HashTable // principle variation table
	PvArray   [MaxDepth]int

	searchHistory [13][BoardSquareNum]int // everytime a search improves alpha, for that piece type and to square, we will improve the score
	searchKillers [2][MaxDepth]int        // stores 2 moves that have recently stored a beta cutoff (not considers captures)
}

// Sq120ToSq64 would return the index of 120 mapped to a 64 square board
var Sq120ToSq64 [BoardSquareNum]int

// Sq64ToSq120 would return the index of 64 mapped to a 120 square board
var Sq64ToSq120 [64]int

// FileRankToSquare converts give file and rank to a square index
func FileRankToSquare(file, rank int) (square int) {
	return ((21 + file) + (rank * 10))
}

// !!!!!!!!! Consider Removing these because the add extra overhead

// Sq64 returns the element at sq120 base
func Sq64(sq120 int) int {
	return Sq120ToSq64[sq120]
}

// Sq120 returns the element at sq64 base
func Sq120(sq64 int) int {
	return Sq64ToSq120[sq64]
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

// SetMask is used when setting a bit to 1 or 0
var SetMask [64]uint64

// ClearMask is used to clear a bit
var ClearMask [64]uint64

// ClearBit takes a bitboard and clears the bit at a provided square
func ClearBit(bb *uint64, sq int) {
	*bb &= ClearMask[sq]
}

// SetBit sets the bit at square to true given a bit board
func SetBit(bb *uint64, sq int) {
	*bb |= SetMask[sq]
}

// PieceKeys hashkeys for each piece for each possible position for the key
var PieceKeys [13][120]uint64

// SideKey the hashkey associated with the current side
var SideKey uint64

// CastleKeys haskeys associated with castling rights
var CastleKeys [16]uint64 // castling value ranges from 0-15 -> we need 16 hashkeys

const (
	// StartFen starting position in fen notation
	StartFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

// PieceNotationMap maps piece notations (i.e. 'p', 'N') to piece values (i.e. 'BlackPawn', 'WhiteKnight')
var PieceNotationMap = map[string]int{
	"p": BlackPawn,
	"r": BlackRook,
	"n": BlackKnight,
	"b": BlackBishop,
	"k": BlackKing,
	"q": BlackQueen,
	"P": WhitePawn,
	"R": WhiteRook,
	"N": WhiteKnight,
	"B": WhiteBishop,
	"K": WhiteKing,
	"Q": WhiteQueen,
}

// FileNotationMap maps file notations (i.e. 'a', 'h') to file values (i.e. 'FileA', 'FileH')
var FileNotationMap = map[string]int{
	"a": FileA,
	"b": FileB,
	"c": FileC,
	"d": FileD,
	"e": FileE,
	"f": FileF,
	"g": FileG,
	"h": FileH,
}

// PieceCharMap maps piece notations (i.e. 'p', 'N') to piece values (i.e. 'BlackPawn', 'WhiteKnight')
var PieceCharMap = map[int]string{
	BlackPawn:   "p",
	BlackRook:   "r",
	BlackKnight: "n",
	BlackBishop: "b",
	BlackKing:   "k",
	BlackQueen:  "q",
	WhitePawn:   "P",
	WhiteRook:   "R",
	WhiteKnight: "N",
	WhiteBishop: "B",
	WhiteKing:   "K",
	WhiteQueen:  "Q",
}

// FilesBoard an array that returns which file a particular square is on
var FilesBoard [BoardSquareNum]int

// RanksBoard an array that returns which file a particular square is on
var RanksBoard [BoardSquareNum]int

// Move type
type Move struct {
	Move  int
	score int
}

/* Game move - information stored in the move int from type Move
   | |-P|-|||Ca-||---To--||-From-|
0000 0000 0000 0000 0000 0111 1111 -> From - 0x7F
0000 0000 0000 0011 1111 1000 0000 -> To - >> 7, 0x7F
0000 0000 0011 1100 0000 0000 0000 -> Captured - >> 14, 0xF
0000 0000 0100 0000 0000 0000 0000 -> En passant capt - 0x40000
0000 0000 1000 0000 0000 0000 0000 -> PawnStart - 0x80000
0000 1111 0000 0000 0000 0000 0000 -> Promotion to what piece - >> 20, 0xF
0001 0000 0000 0000 0000 0000 0000 -> Castle - 0x1000000
*/

// FromSq - macro that returns the 'from' bits from the move int
func FromSq(m int) int {
	return m & 0x7f
}

// ToSq - macro that returns the 'to' bits from the move int
func ToSq(m int) int {
	return (m >> 7) & 0x7f
}

// Captured - macro that returns the 'Captured' bits from the move int
func Captured(m int) int {
	return (m >> 14) & 0xf
}

// Promoted - macro that returns the 'Promoted' bits from the move int
func Promoted(m int) int {
	return (m >> 20) & 0xf
}

const (
	// MoveFlagEnPass move flag that denotes if the capture was an enpass
	MoveFlagEnPass int = 0x40000
	// MoveFlagPawnStart move flag that denotes if move was pawn start (2x)
	MoveFlagPawnStart int = 0x80000
	// MoveFlagCastle move flag that denotes if move was castling
	MoveFlagCastle int = 0x1000000
	// MoveFlagCapture move flag that denotes if move was capture without saying what the capture was (checks capture & enpas squares)
	MoveFlagCapture int = 0x7C000
	// MoveFlagPromotion move flag that denotes if move was promotion without saying what the promotion was
	MoveFlagPromotion int = 0xF00000
)

const (
	// MaxPositionMoves maximum number of posible moves for a given position
	MaxPositionMoves int = 256
)

// MoveList a structure to hold all generated moves
type MoveList struct {
	Moves [MaxPositionMoves]Move
	Count int // number of moves on the moves list
}

// Debug variable that enables/disables debugging
var Debug = true

const (
	// NoMove signifies no move
	NoMove int = 0
)

// HashEntry principle variation entry
type HashEntry struct {
	posKey uint64
	move   int
	score  int
	depth  int
	flags  int
}

// HashTable principle variation table
type HashTable struct {
	pTable     []HashEntry // you can make an array instead but this allows for dynamically allocating space as you go along
	numEntries int
	newWrite   int
	overWrite  int
	hit        int
	cut        int
}

// Hash entry flags
const (
	// HFNone hfnone
	HFNone = iota
	// HFAlpha hf alpha
	HFAlpha
	// HFBeta hf beta
	HFBeta
	// HFExact hf exact
	HFExact
)

const (
	// MaxDepth maximum search depth
	MaxDepth int = 64
)

// SearchInfo struct to hold search related information
type SearchInfo struct {
	StartTime time.Time
	StopTime  int
	Depth     int
	depthSet  int
	TimeSet   bool
	movesToGo int
	infinite  bool // if this is true, do not stop search based on time but when the gui sends the stop command

	nodes uint64 // count of all positions that the engine visits in the search tree

	Quit    bool // if interrupt is sent -> quit
	stopped bool

	failHigh      float32 // these will be used to look at move ordering
	failHighFirst float32
	nullCut       int // null move cutoff

	GameMode     int  // see consts below
	PostThinking bool // if true, engine posts its thinking to the gui
}

// Game Modes
const (
	// UciMode mode using the UCI protocol
	UciMode = iota
	// XBoardMode mode using the XBoard protocol
	XBoardMode
	// ConsoleMode mode using the console for input
	ConsoleMode
)

// FileBBMask evaluation masks that help identify passed pawns, open files etc
var FileBBMask [8]uint64

// RankBBMask evaluation masks that help identify passed pawns, open ranks etc
var RankBBMask [8]uint64

// BlackPassedMask black passed pawn mask
var BlackPassedMask [64]uint64

// WhitePassedMask white passed pawn mask
var WhitePassedMask [64]uint64

// IsolatedMask isolated pawn mask
var IsolatedMask [64]uint64

// WhiteDoubledMask isolated pawn mask
var WhiteDoubledMask [64]uint64

// BlackDoubledMask isolated pawn mask
var BlackDoubledMask [64]uint64

// PieceChar string with piece characters
var PieceChar = ".PNBRQKpnbrqk"

// SideChar string with side characters
var SideChar = "wb-"

// RankChar string with rank characters
var RankChar = "12345678"

// FileChar string with file characters
var FileChar = "abcdefgh"

// PieceBig A map used to identify if a piece is considered "Big"
var PieceBig = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: true,
	WhiteBishop: true,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   true,
	BlackPawn:   false,
	BlackKnight: true,
	BlackBishop: true,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   true,
}

// PieceMaj A map used to identify if a piece is considered "Major"
var PieceMaj = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   true,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   true,
}

// PieceMin A map used to identify if a piece is considered "Minor"
var PieceMin = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: true,
	WhiteBishop: true,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: true,
	BlackBishop: true,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   false,
}

// PieceColour A map used to identify a piece's colour
var PieceColour = map[int]int{
	Empty:       Both,
	WhitePawn:   White,
	WhiteKnight: White,
	WhiteBishop: White,
	WhiteRook:   White,
	WhiteQueen:  White,
	WhiteKing:   White,
	BlackPawn:   Black,
	BlackKnight: Black,
	BlackBishop: Black,
	BlackRook:   Black,
	BlackQueen:  Black,
	BlackKing:   Black,
}

// IsPieceKnight holds information if a given piece is a knight
var IsPieceKnight = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: true,
	WhiteBishop: false,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: true,
	BlackBishop: false,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   false,
}

// IsPieceKing holds information if a given piece is a king
var IsPieceKing = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   true,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   true,
}

// IsPieceRookQueen holds information if a given piece is a rook or queen
var IsPieceRookQueen = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   false,
}

// IsPieceBishopQueen holds information if a given piece is a bishop or queen
var IsPieceBishopQueen = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: true,
	WhiteRook:   false,
	WhiteQueen:  true,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: true,
	BlackRook:   false,
	BlackQueen:  true,
	BlackKing:   false,
}

// IsPiecePawn holds information if a given piece is a pawn
var IsPiecePawn = map[int]bool{
	Empty:       false,
	WhitePawn:   true,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   false,
	BlackPawn:   true,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   false,
}

// PieceSlides holds information if a given piece slides
var PieceSlides = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: true,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: true,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   false,
}

// LoopSlidePiece sliding pieces slice used for looping
var LoopSlidePiece = [...]int{WhiteBishop, WhiteRook, WhiteQueen, 0, BlackBishop, BlackRook, BlackQueen, 0}

// LoopSlideIndex sliding pieces index slice to index where
// the white pieces start in the above LoopSlidePiece, and where black
var LoopSlideIndex = [...]int{0, 4}

// LoopNonSlidePiece non-sliding pieces slice used for looping
var LoopNonSlidePiece = [...]int{WhiteKnight, WhiteKing, 0, BlackKnight, BlackKing, 0}

// LoopNonSlideIndex non-sliding pieces index slice to index where
// the white pieces start in the above LoopSlidePiece, and where black
var LoopNonSlideIndex = [...]int{0, 3}

// PiececeDir squares increment for each direction
var PiececeDir = map[int][]int{
	Empty:       {0, 0, 0, 0, 0, 0, 0},
	WhitePawn:   {0, 0, 0, 0, 0, 0, 0},
	WhiteKnight: {-8, -19, -21, -12, 8, 19, 21, 12},
	WhiteBishop: {-9, -11, 11, 9, 0, 0, 0, 0},
	WhiteRook:   {-1, -10, 1, 10, 0, 0, 0, 0},
	WhiteQueen:  {-1, -10, 1, 10, -9, -11, 11, 9},
	WhiteKing:   {-1, -10, 1, 10, -9, -11, 11, 9},
	BlackPawn:   {0, 0, 0, 0, 0, 0, 0},
	BlackKnight: {-8, -19, -21, -12, 8, 19, 21, 12},
	BlackBishop: {-9, -11, 11, 9, 0, 0, 0, 0},
	BlackRook:   {-1, -10, 1, 10, 0, 0, 0, 0},
	BlackQueen:  {-1, -10, 1, 10, -9, -11, 11, 9},
	BlackKing:   {-1, -10, 1, 10, -9, -11, 11, 9},
}

// NumberOfDir number of directions in which each piece can move
var NumberOfDir = map[int]int{
	Empty:       0,
	WhitePawn:   0,
	WhiteKnight: 8,
	WhiteBishop: 4,
	WhiteRook:   4,
	WhiteQueen:  8,
	WhiteKing:   8,
	BlackPawn:   0,
	BlackKnight: 8,
	BlackBishop: 4,
	BlackRook:   4,
	BlackQueen:  8,
	BlackKing:   8,
}
