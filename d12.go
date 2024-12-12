package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
)

type Garden struct {
	Type  byte
	Items []Pos
}

func (self *Garden) Area() int {
	return len(self.Items)
}

func (self *Garden) Circumference() int {
	xC, yC := genCircumference(self.Items)
	return len(xC) + len(yC)
}

func (self *Garden) String() string {
	xC, yC := genCircumference(self.Items)
	return fmt.Sprintf("(%c: %d %d(%d) %v %v)", self.Type, self.Area(), len(xC)+len(yC), CountHLines(xC)+CountVLines(yC)+countIntersections(self.Items)*2, xC, yC)
}

func (self *Pos) Add(x int, y int) Pos {
	return Pos{X: self.X + x, Y: self.Y + y}
}

func genCircumference(items []Pos) (xs, ys []Pos) {
	xs = []Pos{}
	ys = []Pos{}
	seenXIndex := func(p Pos) (int, bool) {
		return slices.BinarySearchFunc(xs, p, ComparePos)
	}
	seenYIndex := func(p Pos) (int, bool) {
		return slices.BinarySearchFunc(ys, p, ComparePos)
	}

	l := 999999999
	t := 999999999
	w := 0
	h := 0

	for _, pos := range items {
		if pos.X < l {
			l = pos.X
		}
		if pos.Y < t {
			t = pos.Y
		}
		if pos.X > w {
			w = pos.X
		}
		if pos.Y > h {
			h = pos.Y
		}
	}

	containsPos := func(p Pos) bool {
		for _, check := range items {
			if p == check {
				return true
			}
		}
		return false
	}

	for x := l; x <= w; x++ {
		for y := t; y <= h; y++ {
			p := Pos{X: x, Y: y}
			if !containsPos(p) {
				continue
			}

			// Check Xs
			i, found := seenXIndex(p)
			if found {
				xs = slices.Delete(xs, i, i+1)
			} else {
				xs = slices.Insert(xs, i, p)
			}

			if i, found = seenXIndex(p.Add(1, 0)); !found {
				xs = slices.Insert(xs, i, p.Add(1, 0))
			}

			// Check Ys
			i, found = seenYIndex(p)
			if found {
				ys = slices.Delete(ys, i, i+1)
			} else {
				ys = slices.Insert(ys, i, p)
			}

			if i, found = seenYIndex(p.Add(0, 1)); !found {
				ys = slices.Insert(ys, i, p.Add(0, 1))
			}
		}
	}

	return
}

func CountVLines(items []Pos) int {
	l := 999999999
	t := 999999999
	w := 0
	h := 0

	for _, pos := range items {
		if pos.X < l {
			l = pos.X
		}
		if pos.Y < t {
			t = pos.Y
		}
		if pos.X > w {
			w = pos.X
		}
		if pos.Y > h {
			h = pos.Y
		}
	}

	containsPos := func(p Pos) bool {
		for _, check := range items {
			if p == check {
				return true
			}
		}
		return false
	}

	total := 0
	for y := t; y <= h; y++ {

		isFence := false
		for x := l; x <= w; x++ {
			p := Pos{X: x, Y: y}
			if !containsPos(p) {
				isFence = false
			} else {
				if !isFence {
					total += 1
				}
				isFence = true
			}

		}
	}

	return total
}

func CountHLines(items []Pos) int {
	l := 999999999
	t := 999999999
	w := 0
	h := 0

	for _, pos := range items {
		if pos.X < l {
			l = pos.X
		}
		if pos.Y < t {
			t = pos.Y
		}
		if pos.X > w {
			w = pos.X
		}
		if pos.Y > h {
			h = pos.Y
		}
	}

	containsPos := func(p Pos) bool {
		for _, check := range items {
			if p == check {
				return true
			}
		}
		return false
	}

	total := 0

	for x := l; x <= w; x++ {
		isFence := false
		for y := t; y <= h; y++ {
			p := Pos{X: x, Y: y}
			if !containsPos(p) {
				isFence = false
			} else {
				if !isFence {
					total += 1
				}
				isFence = true
			}

		}
	}

	return total
}

func countIntersections(points []Pos) int {
	l := 999999999
	t := 999999999
	w := 0
	h := 0

	for _, pos := range points {
		if pos.X < l {
			l = pos.X
		}
		if pos.Y < t {
			t = pos.Y
		}
		if pos.X > w {
			w = pos.X
		}
		if pos.Y > h {
			h = pos.Y
		}
	}

	containsPos := func(p Pos) bool {
		for _, check := range points {
			if p == check {
				return true
			}
		}
		return false
	}

	grid := make([][]byte, h-t+1)
	for y := range grid {
		grid[y] = make([]byte, w-l+1)
		for x := range grid[y] {
			grid[y][x] = 'O'
		}
	}
	for x := l; x <= w; x++ {
		for y := t; y <= h; y++ {
			p := Pos{X: x, Y: y}
			if !containsPos(p) {
				continue
			}
			grid[y-t][x-l] = 'X'
		}
	}

	patterns := [][][]byte{
		{
			[]byte("XO"),
			[]byte("OX"),
		},
		{
			[]byte("OX"),
			[]byte("XO"),
		},
	}

	total := 0
	for _, pattern := range patterns {
		total += count(grid, pattern)
	}

	return total
}

func scanGarden(garden [][]byte, pos Pos) Garden {
	w, h := shape(garden)

	seen := []Pos{pos}
	toLook := []Pos{pos}

	seenIndex := func(p Pos) (int, bool) {
		return slices.BinarySearchFunc(seen, p, ComparePos)
	}
	g := Garden{Type: garden[pos.Y][pos.X]}

	for len(toLook) > 0 {
		next := toLook[len(toLook)-1]
		toLook = slices.Delete(toLook, len(toLook)-1, len(toLook))

		if garden[next.Y][next.X] == g.Type {
			g.Items = append(g.Items, next)
		} else {
			continue
		}

		if next.X+1 < w {
			n := next.Add(1, 0)
			if i, found := seenIndex(n); !found {
				seen = slices.Insert(seen, i, n)
				toLook = append(toLook, n)
			}
		}
		if next.X-1 >= 0 {
			n := next.Add(-1, 0)
			if i, found := seenIndex(n); !found {
				seen = slices.Insert(seen, i, n)
				toLook = append(toLook, n)
			}
		}
		if next.Y+1 < h {
			n := next.Add(0, 1)
			if i, found := seenIndex(n); !found {
				seen = slices.Insert(seen, i, n)
				toLook = append(toLook, n)
			}
		}
		if next.Y-1 >= 0 {
			n := next.Add(0, -1)
			if i, found := seenIndex(n); !found {
				seen = slices.Insert(seen, i, n)
				toLook = append(toLook, n)
			}
		}
	}

	return g
}

func (self *Day) D12() error {

	garden := [][]byte{}

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

		garden = append(garden, line)

	}

	gardens := []Garden{}
	seenPoses := []Pos{}

	seenIndex := func(p Pos) (int, bool) {
		return slices.BinarySearchFunc(seenPoses, p, ComparePos)
	}

	for y, row := range garden {
		for x := range row {
			p := Pos{X: x, Y: y}
			if _, seen := seenIndex(p); seen {
				continue
			}
			g := scanGarden(garden, p)
			for _, pos := range g.Items {
				if i, found := seenIndex(pos); !found {
					seenPoses = slices.Insert(seenPoses, i, pos)
				}
			}
			gardens = append(gardens, g)
		}
	}

	total := 0
	total2 := 0
	for _, g := range gardens {
		Debug(g.String())
		xs, ys := genCircumference(g.Items)
		total += g.Area() * (len(xs) + len(ys))
		total2 += (CountHLines(xs) + CountVLines(ys) + countIntersections(g.Items)*2) * g.Area()
	}

	fmt.Println(total)

	fmt.Println(total2)

	return nil
}
