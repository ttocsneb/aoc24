package main

import (
	"bytes"
	"fmt"
	"io"
)

func (self *Day) D15() error {

	warehouse := [][]rune{}

	const (
		Box      rune = 'O'
		BoxLeft  rune = '['
		BoxRight rune = ']'
		Robot    rune = '@'
		Wall     rune = '#'
		Empty    rune = '.'
		Up       rune = '^'
		Down     rune = 'v'
		Left     rune = '<'
		Right    rune = '>'
	)

	instructions := []rune{}

	readLine := func() (string, error) {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		line = bytes.TrimSpace(line)
		return string(line), nil
	}
	start, err := readLine()
	if err != nil {
		return err
	}
	w := len(start)
	Debug(w)
	warehouse = append(warehouse, []rune(start))

	// Read warehouse
	for {
		line, err := readLine()
		if err != nil {
			return err
		}
		warehouse = append(warehouse, []rune(line))

		if line == start {
			break
		}
	}

	for {
		line, err := readLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		for _, c := range line {
			instructions = append(instructions, c)
		}
	}

	for _, row := range warehouse {
		Debug(string(row))
	}

	// instruction := 0

	at := func(wh [][]rune, p Pos) rune {
		return wh[p.Y][p.X]
	}

	var canMove func([][]rune, Pos, Pos) bool

	canMove = func(wh [][]rune, from, offset Pos) bool {

		if at(wh, from) == Wall {
			// Debug(" Hit Wall")
			return false
		}
		dest := from.AddSelf(offset)
		// Debugf("Moving %v(%c) -> %v(%c)\n", from, at(wh, from), dest, at(wh, dest))
		if at(wh, dest) != Empty {
			if offset.Y != 0 {
				if at(wh, dest) == BoxLeft {
					if !canMove(wh, dest.Add(1, 0), offset) {
						return false
					}
				} else if at(wh, dest) == BoxRight {
					if !canMove(wh, dest.Add(-1, 0), offset) {
						return false
					}
				}
			}
			if !canMove(wh, dest, offset) {
				return false
			}
		}
		return true
	}

	var moveItem func([][]rune, Pos, Pos) bool
	moveItem = func(wh [][]rune, from, offset Pos) bool {
		if !canMove(wh, from, offset) {
			return false
		}
		dest := from.AddSelf(offset)
		// Debugf("Moving %v(%c) -> %v(%c)\n", from, at(wh, from), dest, at(wh, dest))
		if at(wh, dest) != Empty {
			if offset.Y != 0 {
				if at(wh, dest) == BoxLeft {
					if !moveItem(wh, dest.Add(1, 0), offset) {
						return false
					}
				} else if at(wh, dest) == BoxRight {
					if !moveItem(wh, dest.Add(-1, 0), offset) {
						return false
					}
				}
			}
			if !moveItem(wh, dest, offset) {
				return false
			}
		}
		wh[dest.Y][dest.X] = at(wh, from)
		wh[from.Y][from.X] = Empty
		// Debug("Success")

		return true
	}

	wide := make([][]rune, len(warehouse))
	for y, row := range warehouse {
		wide[y] = make([]rune, len(row)*2)
		for x, c := range row {
			if c == Box {
				wide[y][x*2] = BoxLeft
				wide[y][x*2+1] = BoxRight
			} else if c == Robot {
				wide[y][x*2] = Robot
				wide[y][x*2+1] = Empty
			} else {
				wide[y][x*2] = c
				wide[y][x*2+1] = c
			}
		}
	}

	run := func(wh [][]rune) {
		var robot Pos
		for _, instruction := range instructions {
			for y, row := range wh {
				for x, c := range row {
					if c == Robot {
						robot = Pos{X: x, Y: y}
						goto foundRobot
					}
				}
			}
			panic("Could not find robot")
		foundRobot:
			// Debug(string(instruction))
			switch instruction {
			case Up:
				if moveItem(wh, robot, Pos{X: 0, Y: -1}) {
					robot = robot.AddSelf(Pos{X: 0, Y: -1})
				}
			case Down:
				if moveItem(wh, robot, Pos{X: 0, Y: 1}) {
					robot = robot.AddSelf(Pos{X: 0, Y: 1})
				}
			case Left:
				if moveItem(wh, robot, Pos{X: -1, Y: 0}) {
					robot = robot.AddSelf(Pos{X: -1, Y: 0})
				}
			case Right:
				if moveItem(wh, robot, Pos{X: 1, Y: 0}) {
					robot = robot.AddSelf(Pos{X: 1, Y: 0})
				}
			}
			for _, row := range wh {
				Debug(string(row))
			}
		}
	}
	run(warehouse)

	total := 0
	for y, row := range warehouse {
		for x, c := range row {
			if c == Box {
				total += y*100 + x
			}
		}
	}
	fmt.Println(total)

	for _, row := range wide {
		Debug(string(row))
	}
	run(wide)

	total = 0
	for y, row := range wide {
		for x, c := range row {
			if c == BoxLeft {
				total += y*100 + x
			}
		}
	}
	fmt.Println(total)

	return nil
}
