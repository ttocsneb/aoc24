package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type Vector struct {
	Pos Pos
	Vel Pos
}

func (self Pos) Mul(factor int) Pos {
	return Pos{X: self.X * factor, Y: self.Y * factor}
}

func (self Pos) Mod(x, y int) Pos {
	xm := self.X % x
	ym := self.Y % y
	if xm < 0 {
		xm += x
	}
	if ym < 0 {
		ym += y
	}
	return Pos{X: xm, Y: ym}
}

func (self Vector) Simulate(t int, xBound, yBound int) Pos {
	return self.Pos.AddSelf(self.Vel.Mul(t)).Mod(xBound, yBound)
}

func (self *Day) D14() error {

	vectors := []Vector{}

	readLine := func() (string, error) {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		line = bytes.TrimSpace(line)
		return string(line), nil
	}

	for {
		line, err := readLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		parseNum := func(val string) int {
			numStr := []byte{}
			for _, b := range []byte(val) {
				if b >= '0' && b <= '9' || b == '-' {
					numStr = append(numStr, b)
				}
			}

			n, err := strconv.Atoi(string(numStr))
			if err != nil {
				panic(err)
			}

			return n
		}

		parsePos := func(val string) Pos {
			x, y, _ := strings.Cut(val, ",")

			return Pos{X: parseNum(x), Y: parseNum(y)}
		}

		fields := strings.Fields(line)

		vectors = append(vectors, Vector{
			Pos: parsePos(fields[0]),
			Vel: parsePos(fields[1]),
		})
	}

	// w := 101
	// h := 103
	w := 11
	h := 7

	quads := make([]int, 4)
	wh := w / 2
	hh := h / 2

	for _, vec := range vectors {
		pos := vec.Simulate(100, w, h)

		Debug(vec.Pos, pos)
		if pos.X < wh {
			if pos.Y < hh {
				quads[0] += 1
			} else if pos.Y > hh {
				quads[1] += 1
			}
		} else if pos.X > wh {
			if pos.Y < hh {
				quads[2] += 1
			} else if pos.Y > hh {
				quads[3] += 1
			}
		}
	}

	newGrid := func() [][]byte {
		grid := make([][]byte, h)
		for y := range grid {
			grid[y] = make([]byte, w)
			for x := range grid[y] {
				grid[y][x] = '.'
			}
		}
		return grid
	}

	Debug(quads)
	fmt.Println(quads[0] * quads[1] * quads[2] * quads[3])

	overlaps := func(poses []Pos) bool {
		seen := []Pos{}
		for _, pos := range poses {
			if i, found := slices.BinarySearchFunc(seen, pos, ComparePos); !found {
				seen = slices.Insert(seen, i, pos)
			} else {
				return true
			}
		}
		return false
	}

	for i := 1; ; i++ {
		poses := make([]Pos, len(vectors))
		grid := newGrid()
		for j, vec := range vectors {
			poses[j] = vec.Simulate(i, w, h)
			if grid[poses[j].Y][poses[j].X] == '.' {
				grid[poses[j].Y][poses[j].X] = '1'
			} else {
				grid[poses[j].Y][poses[j].X] += 1
			}

		}
		if overlaps(poses) {
			continue
		} else {
			for _, row := range grid {
				Debug(string(row))
			}
			fmt.Println(i)
			break
		}
	}

	return nil
}
