package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readNext(seq string, cmds ...string) (rest string, op string, cmd string) {
	minPos := len(seq)
	for _, keyword := range cmds {
		i := strings.Index(seq, keyword)
		if i != -1 {
			if minPos > i {
				cmd = keyword
				rest = seq[i+len(keyword):]
				minPos = i
			}
		}
	}
	if minPos == -1 {
		return "", "", ""
	}

	var found bool
	if rest, found = strings.CutPrefix(rest, "("); !found {
		return
	}

	isValid := true
	i := strings.IndexFunc(rest, func(r rune) bool {
		if !isValid {
			return false
		}
		if r == ']' {
			isValid = false
			return false
		}
		if r == '>' {
			isValid = false
			return false
		}
		if r == '(' {
			isValid = false
			return false
		}
		if r == ')' {
			return true
		}
		return false
	})
	if i == -1 {
		return
	}
	op = rest[:i]
	rest = rest[i+1:]
	return
}

func (self *Day) D3() error {
	input, err := io.ReadAll(self.Input)
	if err != nil {
		return err
	}

	data := string(input)
	total := 0
	p2total := 0
	do := true
	for len(data) > 0 {
		var op string
		var cmd string
		data, op, cmd = readNext(data, "mul", "don't", "do")
		if cmd == "do" {
			do = true
			continue
		} else if cmd == "don't" {
			do = false
			continue
		}

		if op != "" {
			a, b, _ := strings.Cut(op, ",")

			ai, err := strconv.ParseInt(a, 10, 0)
			if err != nil {
				continue
			}

			bi, err := strconv.ParseInt(b, 10, 0)
			if err != nil {
				continue
			}
			total += int(ai) * int(bi)

			if do {
				p2total += int(ai) * int(bi)
			} else {
			}
		}
	}

	fmt.Println(total)
	fmt.Println(p2total)

	return nil
}
