package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"
)

type GameState = int

const (
	Waiting = iota
	Running
	Snake1Wins
	Snake2Wins
	Draw
)

type Grid = [][]byte

type Cell struct {
	row int
	col int
}

type Direction = int

const (
	Up = iota
	Down
	Left
	Right
)

type Snake = []Cell
type Head = Cell

type Player struct {
	snake Snake
	dir   Direction
	score int
}

func NewPlayer(s Snake, dir Direction, score int) *Player {
	return &Player{
		snake: s,
		dir:   dir,
		score: score,
	}
}

func (p *Player) dead(other Snake, rows, cols int) bool {
	h := head(p.snake)
	if h.row < 1 || h.row >= rows-1 || h.col < 1 || h.col >= cols-1 {
		return true
	}

	if slices.Contains(headless(p.snake), h) || slices.Contains(other, h) {
		return true
	}
	return false
}

func (p *Player) eaten(f Fruit) bool {
	eaten := head(p.snake) == f
	if eaten {
		p.score += 1
	}
	return eaten
}

func (p *Player) move(eaten bool) {
	newHead := head(p.snake)
	switch p.dir {
	case Up:
		newHead.row -= 1
	case Down:
		newHead.row += 1
	case Left:
		newHead.col -= 1
	case Right:
		newHead.col += 1
	}
	if eaten {
		p.snake = append(p.snake, newHead)
	} else {
		l := len(p.snake) - 1
		for pos := range len(p.snake) - 1 {
			p.snake[pos] = p.snake[pos+1]
		}
		p.snake[l] = newHead
	}
}

func head(s Snake) Head {
	return s[len(s)-1]
}

func headless(s Snake) Snake {
	return s[:len(s)-1]
}

type Fruit = Cell

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

func drawFruit(grid Grid, f Fruit) {
	grid[f.row][f.col] = '*'
}

func drawSnake(grid Grid, s Snake) {
	h := head(s)
	for _, c := range s {
		grid[c.row][c.col] = '#'
	}
	grid[h.row][h.col] = 'o'
}

func render(grid [][]byte, p1, p2 *Player, fruit Fruit, state GameState) {
	rows, cols := len(grid)-1, len(grid[0])-1
	reset(grid, rows, cols)

	drawBorder(grid, rows, cols)
	drawFruit(grid, fruit)
	drawSnake(grid, p1.snake)
	drawSnake(grid, p2.snake)

	fmt.Printf("score 1: %d | score 2: %d\r\n", p1.score, p2.score)
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Printf("\r\n")
	}

	if state != Running {
		switch state {
		case Draw:
			fmt.Printf("\033[%dA", rows/2+3)
			fmt.Printf("\033[%dC", cols/2-1)
			fmt.Printf("Draw\r\n")
		case Snake1Wins:
			fmt.Printf("\033[%dA", rows/2+3)
			fmt.Printf("\033[%dC", cols/2-4)
			fmt.Printf("Player 1 Wins\r\n")
		case Snake2Wins:
			fmt.Printf("\033[%dA", rows/2+3)
			fmt.Printf("\033[%dC", cols/2-4)
			fmt.Printf("Player 2 Wins\r\n")
		}
		fmt.Printf("\033[%dC", cols/2-7)
		fmt.Printf("Press R to restart\r\n")

		fmt.Printf("\033[%dB", rows/2+1)
		fmt.Printf("\033[%dD", cols)
	}
	fmt.Printf("\033[%dA", rows+2)
}

func newFruit(rows, cols int, snake1, snake2 Snake) Fruit {
	row := rand.Intn(rows-2) + 1 // [1, rows-2]
	col := rand.Intn(cols-2) + 1 // [1, cols-2]
	// TODO: could be bad when game is close to finishing
	for _, s := range snake1 {
		if s.row == row && s.col == col {
			return newFruit(rows, cols, snake1, snake2)
		}
	}
	for _, s := range snake2 {
		if s.row == row && s.col == col {
			return newFruit(rows, cols, snake1, snake2)
		}
	}
	return Cell{row, col}
}

func runGame() {
	rows, cols := 30, 60
	grid := make([][]byte, rows)
	for row := range grid {
		grid[row] = make([]byte, cols)
	}

	state := Waiting

	var p1 *Player
	var p2 *Player
	var fruit Fruit

	startGame := func() {
		s1 := Snake{{10, 23}, {10, 22}, {10, 21}, {10, 20}}
		s2 := Snake{{13, 30}, {13, 29}, {10, 28}, {10, 27}}
		p1 = NewPlayer(s1, Left, 0)
		p2 = NewPlayer(s2, Up, 0)
		fruit = newFruit(rows, cols, p1.snake, p2.snake)
		state = Running
	}
	startGame()

	// go func() {
	// 	for {
	// 		buf := make([]byte, 1)
	// 		n, err := os.Stdin.Read(buf)
	// 		if err != nil {
	// 			cleanup()
	// 			log.Fatal(err)
	// 		}
	// 		if n > 0 {
	// 			switch buf[0] {
	// 			case 3: // ^C
	// 				cleanup()
	// 				os.Exit(0)

	// 			// TODO: fix bug where the player switches directions quickly
	// 			// and the snake moves backwards
	// 			case 'a':
	// 				if p1.dir != Right {
	// 					p1.dir = Left
	// 				}
	// 			case 'w':
	// 				if p1.dir != Down {
	// 					p1.dir = Up
	// 				}
	// 			case 'd':
	// 				if p1.dir != Left {
	// 					p1.dir = Right
	// 				}
	// 			case 's':
	// 				if p1.dir != Up {
	// 					p1.dir = Down
	// 				}

	// 			case 'j':
	// 				if p2.dir != Right {
	// 					p2.dir = Left
	// 				}
	// 			case 'i':
	// 				if p2.dir != Down {
	// 					p2.dir = Up
	// 				}
	// 			case 'l':
	// 				if p2.dir != Left {
	// 					p2.dir = Right
	// 				}
	// 			case 'k':
	// 				if p2.dir != Up {
	// 					p2.dir = Down
	// 				}

	// 			case 'r':
	// 				if state != Running {
	// 					startGame()
	// 				}
	// 			}
	// 		}
	// 	}
	// }()

	for {
		snake1Dead := p1.dead(p2.snake, rows, cols)
		snake2Dead := p2.dead(p1.snake, rows, cols)

		if snake1Dead && snake2Dead {
			if p1.score == p2.score {
				state = Draw
			} else if p1.score > p2.score {
				state = Snake1Wins
			} else {
				state = Snake2Wins
			}
		}

		if state == Running {
			eaten1 := p1.eaten(fruit)
			eaten2 := p2.eaten(fruit)
			if !snake1Dead {
				p1.move(eaten1)
			}
			if !snake2Dead {
				p2.move(eaten2)
			}
			if eaten1 || eaten2 {
				fruit = newFruit(rows, cols, p1.snake, p2.snake)
			}
		}

		render(grid, p1, p2, fruit, state)
		time.Sleep(100 * time.Millisecond)
	}
}
