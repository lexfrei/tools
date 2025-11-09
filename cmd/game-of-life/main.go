package main

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

//nolint:forbidigo,mnd // this file is a draft for the game of life
func main() {
	fmt.Println("Game of Life")

	field := createField(57, 242)

	generation := 1

	var snapshots map[int][][]bool

	snapshots = make(map[int][][]bool)

	for {
		fmt.Print(printField(field))

		fmt.Println("Generation:", generation)

		generation++

		if generation%5 == 0 || generation%3 == 0 {
			snapshots = saveSnapshot(field, snapshots)
		}

		field = nextGeneration(field)

		time.Sleep(100 * time.Millisecond)
		clearScreen(field)
	}
}

// createField creates a field of size x size
// randomly populating it with alive cells (30% chance).
func createField(width, height int) [][]bool {
	field := make([][]bool, width)

	for i := range field {
		field[i] = make([]bool, height)
		for j := range field[i] {
			field[i][j] = randBool()
		}
	}

	return field
}

// printField prints the field to the textView.
func printField(field [][]bool) string {
	var result string

	var resultSb59 strings.Builder
	for i := range field {
		var resultSb60 strings.Builder
		for j := range field[i] {
			if field[i][j] {
				resultSb60.WriteString("X")
			} else {
				resultSb60.WriteString(" ")
			}
		}
		result += resultSb60.String()

		resultSb59.WriteString("\n")
	}
	result += resultSb59.String()

	return result
}

func randBool() bool {
	dataBytes := make([]byte, 1)

	_, err := rand.Read(dataBytes)
	if err != nil {
		slog.Info("error generating random bool", "error", err)

		return false
	}

	return dataBytes[0]%3 == 0
}

func countNeighbours(field [][]bool, x, y int) int {
	var count int

	for counter := x - 1; counter <= x+1; counter++ {
		for j := y - 1; j <= y+1; j++ {
			if (counter == x && j == y) || counter < 0 || j < 0 || counter >= len(field) || j >= len(field[counter]) {
				continue
			}

			if field[counter][j] {
				count++
			}
		}
	}

	return count
}

//nolint:mnd // 3 is not a magic number
func nextGeneration(field [][]bool) [][]bool {
	next := make([][]bool, len(field))

	for row := range field {
		next[row] = make([]bool, len(field[row]))

		for column := range field[row] {
			count := countNeighbours(field, row, column)

			if field[row][column] {
				next[row][column] = count == 2 || count == 3
			} else {
				next[row][column] = count == 3
			}
		}
	}

	return next
}

// clearScreen clears the terminal screen
//
//nolint:forbidigo,goconst // this function is used for clearing the screen
func clearScreen(field [][]bool) {
	clearString := "\033[1A\033[2K"

	// repeat for every line
	var clearStringSb133 strings.Builder
	for range field {
		clearStringSb133.WriteString("\033[1A\033[2K")
	}
	clearString += clearStringSb133.String()

	fmt.Print(clearString) // move cursor up and clear line
}

func saveSnapshot(field [][]bool, snapshots map[int][][]bool) map[int][][]bool {
	snapshots[len(snapshots)] = field

	return snapshots
}
