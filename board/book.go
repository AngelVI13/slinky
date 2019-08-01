package board

import (
	"bufio"
	"log"
	"os"
	"fmt"
	"local/string-utils"
	"math/rand"
	"strings"
	"time"
)

// ScanFile reads file and returns []slice with all lines
func ScanFile(filename string) ([]string, error) {
	lines := make([]string, 0)

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Printf("open file error: %v", err)
		return lines, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return nil, err
	}

	return lines, nil
}


// GetBookMove returns a move from the opening book
func GetBookMove(pos *ChessBoard) int {
	if pos.histPly > 25 {
		return 0
	}

	book, err := ScanFile(BookFile)
	if err != nil {
		fmt.Println("Book error")
		return 0
	}

	currentLine := ""
	for i := 0; i < pos.histPly; i++ {
		currentLine += PrintMove(pos.history[i].move) + " "
	}

	fmt.Println(currentLine)
	bookMoves := make([]int, 0)

	for _, bookLine := range book {
		if strings.Contains(bookLine, currentLine) {
			nextMovesStr := stringutils.RemoveStringToTheLeftOfMarker(bookLine, currentLine)
			nextMoveStr := stringutils.RemoveStringToTheRightOfMarker(nextMovesStr, " ")

			if len(nextMoveStr) > 5 {
				fmt.Println("Book move parsing error")
				continue // parsing eror
			} else {
				bookMoves = append(bookMoves, pos.ParseMove(nextMoveStr))
			}
		}
	}

	numberOfBookMoves := len(bookMoves)
	if len(bookMoves) > 0 {
		rand.Seed(time.Now().Unix())
		return bookMoves[rand.Intn(numberOfBookMoves)]
	}
	return 0
}
