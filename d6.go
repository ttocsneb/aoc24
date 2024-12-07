package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
)

type Pos struct {
	X int
	Y int
}

type Dir int

const (
	Up Dir = iota
	Right
	Down
	Left
)

func ComparePos(lhs, rhs Pos) int {
	if lhs.X < rhs.X {
		return 1
	}
	if lhs.X > rhs.X {
		return -1
	}
	if lhs.Y < rhs.Y {
		return 1
	}
	if lhs.Y > rhs.Y {
		return -1
	}
	return 0
}

func maxDims(blocks []Pos) (x, y int) {
	for _, block := range blocks {
		if block.Y > y {
			y = block.Y
		}
		if block.X > x {
			x = block.X
		}
	}
	return
}

func printLayout(blocks []Pos, guard Pos, direction Dir, visited []Pos) {
	y, x := maxDims(blocks)
	y2, x2 := maxDims(visited)
	x = max(x, x2, guard.X)
	y = max(y, y2, guard.Y)

	buf := make([][]byte, y+1)
	for y := range buf {
		buf[y] = make([]byte, x+1)
		for x := range buf[y] {
			buf[y][x] = '.'
		}
	}

	for _, block := range blocks {
		buf[block.Y][block.X] = '#'
	}
	for _, block := range visited {
		if buf[block.Y][block.X] != '.' {
			buf[block.Y][block.X] = '!'
			continue
		}
		buf[block.Y][block.X] = 'X'
	}

	var g byte
	switch direction {
	case Up:
		g = '^'
	case Right:
		g = '>'
	case Down:
		g = 'v'
	case Left:
		g = '<'
	}
	if buf[guard.Y][guard.X] == '#' {
		buf[guard.Y][guard.X] = '@'
	} else {

		buf[guard.Y][guard.X] = g
	}

	for _, line := range buf {
		fmt.Println(string(line))
	}
}

func AddPos(poses []Pos, pos Pos) []Pos {
	i, found := slices.BinarySearchFunc(poses, pos, ComparePos)

	if !found {
		return slices.Insert(poses, i, pos)
	}

	return poses
}

func HasPos(poses []Pos, pos Pos) bool {
	_, found := slices.BinarySearchFunc(poses, pos, ComparePos)
	return found
}

func CountPositions(blocks []Pos, guard Pos) ([]Pos, bool) {
	getNextPos := func(pos Pos, dir Dir) Pos {
		nextPos := pos
		switch dir {
		case Up:
			nextPos.Y -= 1
		case Right:
			nextPos.X += 1
		case Down:
			nextPos.Y += 1
		case Left:
			nextPos.X -= 1
		default:
			panic("impossible direction")
		}
		return nextPos
	}
	// printLayout(blocks, guard, Up, []Pos{})

	h, w := maxDims(blocks)
	w = max(w, guard.X)
	h = max(h, guard.Y)

	visited := []Pos{guard}
	direction := Up
	curPos := guard
	iter := 0
	lastN := len(visited)
	for {
		iter += 1
		if iter > len(visited)*2 {
			// printLayout(blocks, curPos, direction, visited)
			return visited, true
		}
		if len(visited) != lastN {
			iter = 0
			lastN = len(visited)
		}
		// fmt.Println()
		nextPos := getNextPos(curPos, direction)
		if HasPos(blocks, nextPos) {
			switch direction {
			case Up:
				direction = Right
			case Right:
				direction = Down
			case Down:
				direction = Left
			case Left:
				direction = Up
			default:
				panic("impossible direction")
			}
		} else {
			curPos = nextPos
			if curPos.X <= h && curPos.X >= 0 &&
				curPos.Y <= w && curPos.Y >= 0 {
				visited = AddPos(visited, curPos)
			} else {
				return visited, false
			}
		}

		if curPos.X == guard.X && curPos.Y == guard.Y && direction == Up {
			// printLayout(blocks, curPos, direction, visited)
			return visited, true
		}
	}
}

func (self *Day) D6() error {
	// total := 0
	blocks := []Pos{}
	guard := Pos{}

	h := 0
	for {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			break
		}

		for x, c := range line {
			if c == '.' {
				continue
			}
			if c == '#' {
				blocks = AddPos(blocks, Pos{
					X: x,
					Y: h,
				})
			} else if c == '^' {
				guard = Pos{X: x, Y: h}
			}
		}
		h += 1
	}

	visited, _ := CountPositions(blocks, guard)
	fmt.Println(len(visited))

	h, w := maxDims(blocks)
	w = max(w, guard.X)
	h = max(h, guard.Y)

	count := 0
	lastPct := 0
	fmt.Print("0%")
	for i, pos := range visited {
		if pct := 100 * i / len(visited); pct != lastPct {
			lastPct = pct
			fmt.Printf("\033[2K\r%d%%", lastPct)
		}
		// fmt.Println(pos)
		altered := make([]Pos, len(blocks))
		for i, b := range blocks {
			altered[i] = b
		}
		altered = AddPos(altered, pos)
		if len(altered) == len(blocks) {
			continue
		}
		_, loops := CountPositions(altered, guard)
		if loops {
			count += 1
		}
	}
	fmt.Printf("\033[2K\r%d\n", count)

	return nil
}
