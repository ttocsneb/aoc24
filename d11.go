package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func processRule(m map[int]int, val int, count int) {
	if val == 0 {
		m[1] = m[1] + count
		return
	}

	str := strconv.Itoa(val)
	if len(str)&1 == 0 {
		a, _ := strconv.Atoi(str[:len(str)/2])
		b, _ := strconv.Atoi(str[len(str)/2:])
		m[a] = m[a] + count
		m[b] = m[b] + count
		return
	}

	m[val*2024] += count
}

var seenValues map[int]map[int]int

func (self *Day) D11() error {

	stones := map[int]int{}

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

		for _, field := range strings.Fields(string(line)) {
			val, err := strconv.Atoi(field)
			if err != nil {
				return err
			}
			stones[val] = stones[val] + 1
		}
	}

	Debug(stones)

	countStones := func() int {
		total := 0
		for _, count := range stones {
			total += count
		}
		return total
	}

	for i := 0; i < 75; i++ {
		if i == 25 {
			fmt.Println(countStones())
		}
		nextState := map[int]int{}
		for stone, count := range stones {
			if count == 0 {
				continue
			}
			processRule(nextState, stone, count)
		}
		stones = nextState
		Debugf("%v %v\n", countStones(), stones)
	}

	fmt.Println(countStones())

	return nil
}
