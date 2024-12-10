package main

import (
	"bytes"
	"fmt"
	"io"
)

func arr(n int) []byte {
	a := make([]byte, n)
	for i := range a {
		a[i] = ' '
	}
	return a
}

func print(data [][]byte) {
	if !isDebug {
		return
	}
	for _, row := range data {
		for _, c := range row {
			if c > 0x20 && c < 0x7F {
				fmt.Printf("%c", c)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func compare(block [][]byte, pattern [][]byte) bool {
	for y := range pattern {
		for x := range pattern[y] {
			if pattern[y][x] != '.' {
				if block[y][x] != pattern[y][x] {
					return false
				}
			}
		}
	}
	return true
}

func look(table [][]byte, x, y int, template [][]byte) {
	for i := range template {
		for j := range template[i] {
			template[i][j] = table[i+y][j+x]
		}
	}
}
func shape[T any](data [][]T) (w, h int) {
	return len(data[0]), len(data)
}

func count(table [][]byte, pattern [][]byte, projection ...[][]byte) int {
	w, h := shape(pattern)
	total := 0

	tw, th := shape(table)

	template := make([][]byte, h)
	for y := range template {
		template[y] = make([]byte, w)
	}

	for x := 0; x <= tw-w; x++ {
		for y := 0; y <= th-h; y++ {
			look(table, x, y, template)
			if compare(template, pattern) {
				total += 1
				if len(projection) > 0 {
					project(projection[0], pattern, x, y)
					if isDebug {
						print(projection[0])
					}
				}
			}
		}
	}

	return total
}

func project(table [][]byte, pattern [][]byte, x, y int) {
	for i, row := range pattern {
		for j, c := range row {
			if c != '.' {
				table[y+i][x+j] = c
			}
		}
	}
}

func (self *Day) D4() error {

	total := 0
	table := [][]byte{}
	for {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		table = append(table, line)
		line = bytes.TrimSpace(line)
	}

	// w, h := shape(table)
	//
	// projection := make([][]byte, h)
	// for y := range projection {
	// 	projection[y] = arr(w)
	// }

	XmasPatterns := [][][]byte{
		{
			[]byte("XMAS"),
		},
		{
			[]byte("SAMX"),
		},
		{
			{'X'},
			{'M'},
			{'A'},
			{'S'},
		},
		{
			{'S'},
			{'A'},
			{'M'},
			{'X'},
		},
		{
			[]byte("S..."),
			[]byte(".A.."),
			[]byte("..M."),
			[]byte("...X"),
		},
		{
			[]byte("X..."),
			[]byte(".M.."),
			[]byte("..A."),
			[]byte("...S"),
		},
		{
			[]byte("...X"),
			[]byte("..M."),
			[]byte(".A.."),
			[]byte("S..."),
		},
		{
			[]byte("...S"),
			[]byte("..A."),
			[]byte(".M.."),
			[]byte("X..."),
		},
	}

	for _, pattern := range XmasPatterns {
		total += count(table, pattern)
	}

	if isDebug {
		print(table)
	}

	fmt.Println(total)

	MasPatterns := [][][]byte{
		{
			[]byte("M.S"),
			[]byte(".A."),
			[]byte("M.S"),
		},
		{
			[]byte("S.S"),
			[]byte(".A."),
			[]byte("M.M"),
		},
		{
			[]byte("S.M"),
			[]byte(".A."),
			[]byte("S.M"),
		},
		{
			[]byte("M.M"),
			[]byte(".A."),
			[]byte("S.S"),
		},
	}
	masTotal := 0
	for _, pattern := range MasPatterns {
		masTotal += count(table, pattern)
	}

	fmt.Println(masTotal)
	return nil
}
