package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
)

type Finder struct {
	Grid [][]int
	Seen []Pos
	Pos  Pos
}

func (self *Finder) At(pos Pos) int {
	return self.Grid[pos.Y][pos.X]
}
func (self *Finder) AtXY(x int, y int) int {
	return self.Grid[y][x]
}

func (self *Finder) NextLocs() []Pos {
	x := self.Pos.X
	y := self.Pos.Y
	w, h := shape(self.Grid)

	alt := self.At(self.Pos)

	check := []Pos{}
	if x-1 >= 0 {
		if self.AtXY(x-1, y)-alt == 1 {
			check = append(check, Pos{X: x - 1, Y: y})
		}
	}
	if x+1 < w {
		if self.AtXY(x+1, y)-alt == 1 {
			check = append(check, Pos{X: x + 1, Y: y})
		}
	}
	if y-1 >= 0 {
		if self.AtXY(x, y-1)-alt == 1 {
			check = append(check, Pos{X: x, Y: y - 1})
		}
	}
	if y+1 < h {
		if self.AtXY(x, y+1)-alt == 1 {
			check = append(check, Pos{X: x, Y: y + 1})
		}
	}

	return check
}

func (self *Finder) Pathfind() [][]Pos {
	branches := [][]Pos{}
	for _, possible := range self.NextLocs() {
		if slices.Contains(self.Seen, possible) {
			continue
		}

		cpy := make([]Pos, len(self.Seen), len(self.Seen)+1)
		for i, pos := range self.Seen {
			cpy[i] = pos
		}
		cpy = append(cpy, possible)
		if self.At(possible) == 9 {
			branches = append(branches, cpy)
			continue
		}

		branch := Finder{
			Grid: self.Grid,
			Seen: append(cpy, possible),
			Pos:  possible,
		}
		branches = append(branches, branch.Pathfind()...)
	}
	return branches
}

func (self *Day) D10() error {

	grid := [][]int{}

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

		row := make([]int, len(line))

		for i, c := range line {
			row[i] = int(c - '0')
		}

		grid = append(grid, row)
	}

	startPoses := []Pos{}

	for y, row := range grid {
		for x, alt := range row {
			if alt == 0 {
				startPoses = append(startPoses, Pos{X: x, Y: y})
			}
		}
	}

	total := 0
	total2 := 0
	for _, start := range startPoses {
		trail := Finder{
			Grid: grid,
			Seen: []Pos{},
			Pos:  start,
		}
		paths := trail.Pathfind()
		poses := []Pos{}
		for _, path := range paths {
			last := path[len(path)-1]
			if i, found := slices.BinarySearchFunc(poses, last, ComparePos); !found {
				poses = slices.Insert(poses, i, last)
			}
		}
		total2 += len(paths)
		total += len(poses)
		Debug(start, ":", len(poses))
		Debug(" ", poses)
	}

	fmt.Println(total)
	fmt.Println(total2)

	return nil
}
