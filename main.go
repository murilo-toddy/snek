package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

type Grid = [][]byte

type Cell struct {
	row int
	col int
}

type Snake = []Cell

func head(snake Snake) Cell {
	return snake[len(snake)-1]
}

func headless(snake Snake) Snake {
	return snake[:len(snake)-1]
}

type Fruit = Cell

type Direction = int

const (
	Up = iota
	Down
	Left
	Right
)

func assert(condition bool, msg string) {
	if !condition {
		panic(msg)
	}
}

func drawBorder(grid Grid, rows, cols int) {
	grid[0][0] = '/'
	grid[0][cols] = '\\'
	grid[rows][cols] = '/'
	grid[rows][0] = '\\'

	for row := range rows - 1 {
		grid[row+1][0] = '|'
		grid[row+1][cols] = '|'
	}
	for col := range cols - 1 {
		grid[0][col+1] = '-'
		grid[rows][col+1] = '-'
	}
}

func reset(grid Grid, rows, cols int) {
	for row := range rows {
		for col := range cols {
			grid[row][col] = ' '
		}
	}
}

func drawFruit(grid Grid, pos Cell) {
	grid[pos.row][pos.col] = '*'
}

func drawSnake(snake Snake, grid Grid) {
	h := head(snake)
	for _, cell := range snake {
		grid[cell.row][cell.col] = '#'
	}
	grid[h.row][h.col] = 'o'
}

func render(grid [][]byte, snake Snake, fruit Fruit, score int, lost bool) {
	assert(len(grid) > 0, "grid should not have 0 rows")
	rows, cols := len(grid)-1, len(grid[0])-1
	reset(grid, rows, cols)

	drawBorder(grid, rows, cols)
	drawFruit(grid, fruit)
	drawSnake(snake, grid)

	fmt.Printf("score: %d\r\n", score)
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Printf("\r\n")
	}
	if lost {
		fmt.Printf("\033[%dA", rows/2+3)
		fmt.Printf("\033[%dC", cols/2-3)
		fmt.Printf("You lost\r\n")

		fmt.Printf("\033[%dC", cols/2-7)
		fmt.Printf("Press R to restart\r\n")

		fmt.Printf("\033[%dB", rows/2+1)
		fmt.Printf("\033[%dD", cols)
	}
	fmt.Printf("\033[%dA", rows+2)
}

func move(snake *Snake, dir Direction, fruitEaten bool) {
	newHead := head(*snake)
	switch dir {
	case Up:
		newHead.row -= 1
	case Down:
		newHead.row += 1
	case Left:
		newHead.col -= 1
	case Right:
		newHead.col += 1
	}
	if fruitEaten {
		*snake = append(*snake, newHead)
	} else {
		l := len(*snake) - 1
		for pos := range len(*snake) - 1 {
			(*snake)[pos] = (*snake)[pos+1]
		}
		(*snake)[l] = newHead
	}
}

func newFruit(rows, cols int, snake Snake) Fruit {
	row := rand.Intn(rows-2) + 1 // [1, rows-2]
	col := rand.Intn(cols-2) + 1 // [1, cols-2]
	// TODO: could be bad when game is close to finishing
	for _, s := range snake {
		if s.row == row && s.col == col {
			return newFruit(rows, cols, snake)
		}
	}
	return Cell{row, col}
}

func dead(snake Snake, rows, cols int) bool {
	h := head(snake)
	if h.row < 1 || h.row >= rows-1 || h.col < 1 || h.col >= cols-1 {
		return true
	}
	for _, b := range headless(snake) {
		if b == h {
			return true
		}
	}
	return false
}

func cleanupFunc(oldState *term.State) func() {
	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Printf("\033[?25h") // restore cursor
	}
}

func main() {
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		log.Fatal(fmt.Errorf("stdin/stdout should be terminal"))
	}
	fmt.Printf("\033[?25l") // remove cursor

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	cleanup := cleanupFunc(oldState)
	if err != nil {
		cleanup()
		log.Fatal(err)
	}

	rows, cols := 20, 40
	grid := make([][]byte, rows)
	for row := range grid {
		grid[row] = make([]byte, cols)
	}

	lost := false

	var dir Direction
	var snake Snake
	var fruit Fruit
	var score int

	setup := func() {
		dir = Left
		snake = Snake{{10, 23}, {10, 22}, {10, 21}, {10, 20}}
		fruit = newFruit(rows, cols, snake)
		score = 0
	}
	setup()

	go func() {
		for {
			buf := make([]byte, 1)
			n, err := os.Stdin.Read(buf)
			if err != nil {
				cleanup()
				log.Fatal(err)
			}
			if n > 0 {
				switch buf[0] {
				case 3: // ^C
					cleanup()
					os.Exit(0)

				case 'a':
					if dir != Right {
						dir = Left
					}
				case 'w':
					if dir != Down {
						dir = Up
					}
				case 'd':
					if dir != Left {
						dir = Right
					}
				case 's':
					if dir != Up {
						dir = Down
					}
				case 'r':
					if lost {
						setup()
						lost = false
					}
				}
			}
		}
	}()

	for {
		lost = dead(snake, rows, cols)

		if !lost {
			fruitEaten := head(snake) == fruit
			move(&snake, dir, fruitEaten)
			if fruitEaten {
				fruit = newFruit(rows, cols, snake)
				score += 1
			}
		}

		render(grid, snake, fruit, score, lost)
		time.Sleep(100 * time.Millisecond)
	}
}
