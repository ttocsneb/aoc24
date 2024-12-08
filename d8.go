package main

import (
	"bytes"
	"fmt"
	"io"
	"slices"
)

func printNodes(nodes map[byte][]Pos, antinodes []Pos) {

	x, y := maxDims(antinodes)
	xm, ym := minDims(antinodes)

	for _, nds := range nodes {
		x2, y2 := maxDims(nds)
		xm2, ym2 := minDims(nds)
		y = max(y, y2)
		x = max(x, x2)
		ym = min(ym, ym2)
		xm = min(xm, xm2)
	}

	buf := make([][]byte, y-ym+1)
	for y := range buf {
		buf[y] = make([]byte, x-xm+1)
		for x := range buf[y] {
			buf[y][x] = '.'
		}
	}

	for _, block := range antinodes {
		buf[block.Y-ym][block.X-xm] = '#'
	}
	for k, ps := range nodes {
		for _, pos := range ps {
			if buf[pos.Y-ym][pos.X-xm] != '.' {
				buf[pos.Y-ym][pos.X-xm] = '!'
				continue
			}
			buf[pos.Y-ym][pos.X-xm] = k
		}
	}

	for _, line := range buf {
		fmt.Println(string(line))
	}
}

func genAntinodes(nodes []Pos) []Pos {
	antinodes := []Pos{}

	for a := 0; a < len(nodes); a++ {
		for b := 0; b < len(nodes); b++ {
			if a == b {
				continue
			}
			na := nodes[a]
			nb := nodes[b]
			x := na.X - nb.X
			y := na.Y - nb.Y

			antinodes = append(antinodes, Pos{
				X: na.X + x,
				Y: na.Y + y,
			})
		}
	}

	return antinodes
}

func genAntinodesExtended(nodes []Pos, w, h int) []Pos {
	antinodes := []Pos{}

	for a := 0; a < len(nodes); a++ {
		for b := 0; b < len(nodes); b++ {
			if a == b {
				continue
			}
			na := nodes[a]
			nb := nodes[b]
			x := na.X
			y := na.Y
			xd := na.X - nb.X
			yd := na.Y - nb.Y

			for x < w && x >= -1 && y < h && y >= 0 {
				antinodes = append(antinodes, Pos{
					X: x,
					Y: y,
				})
				x += xd
				y += yd
			}

		}
	}

	return antinodes
}

func (self *Day) D8() error {
	nodes := map[byte][]Pos{}

	h := 0
	w := 0
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
		w = len(line)

		for x, c := range line {
			if c == '.' {
				continue
			}

			n, exists := nodes[c]
			if !exists {
				n = []Pos{}
			}
			n = append(n, Pos{X: x, Y: h})
			nodes[c] = n
		}
		h += 1
	}

	antinodes := []Pos{}
	for _, nds := range nodes {
		antinodes = append(antinodes, genAntinodes(nds)...)
	}

	actuals := []Pos{}
	for _, node := range antinodes {
		if node.X < 0 || node.Y < 0 || node.X >= w || node.Y >= h {
			continue
		}
		i, found := slices.BinarySearchFunc(actuals, node, ComparePos)
		if !found {
			actuals = slices.Insert(actuals, i, node)
		}
	}

	antinodes = []Pos{}
	for _, nds := range nodes {
		antinodes = append(antinodes, genAntinodesExtended(nds, w, h)...)
	}

	actuals2 := []Pos{}
	for _, node := range antinodes {
		if node.X < 0 || node.Y < 0 || node.X >= w || node.Y >= h {
			continue
		}
		i, found := slices.BinarySearchFunc(actuals2, node, ComparePos)
		if !found {
			actuals2 = slices.Insert(actuals2, i, node)
		}
	}

	// printNodes(nodes, []Pos{})
	// fmt.Println()
	// printNodes(nodes, actuals)
	// fmt.Println()
	// printNodes(nodes, actuals2)

	fmt.Println(len(actuals))
	fmt.Println(len(actuals2))

	return nil
}
