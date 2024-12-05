package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
)

type Pair struct {
	Before int
	After  int
}

type PageSet struct {
	Pages []int
	Rules []Pair
}

func (self *PageSet) IsOrdered() bool {
	for _, rule := range self.Rules {
		if !rule.InOrder(self.Pages) {
			return false
		}
	}
	return true
}

func (self *PageSet) Copy() PageSet {
	cpy := PageSet{
		Pages: make([]int, len(self.Pages)),
		Rules: self.Rules,
	}
	for i, val := range self.Pages {
		cpy.Pages[i] = val
	}
	return cpy
}

func (self *PageSet) Order() {
	i := 0
	for !self.IsOrdered() {
		if i >= len(self.Rules) {
			i = 0
		}
		for ; i < len(self.Rules); i++ {
			rule := self.Rules[i]
			if rule.InOrder(self.Pages) {
				continue
			}

			before, after := rule.Indexes(self.Pages)
			self.Pages[before], self.Pages[after] = self.Pages[after], self.Pages[before]
			break
		}
	}
}

func (self *Pair) ContainsBoth(list []int) bool {
	if !slices.Contains(list, self.Before) {
		return false
	}
	if !slices.Contains(list, self.After) {
		return false
	}
	return true
}

func (self *Pair) InOrder(list []int) bool {
	before := slices.Index(list, self.Before)
	if before == -1 {
		return true
	}
	after := slices.Index(list, self.After)
	if after == -1 {
		return true
	}

	return before < after
}

func (self *Pair) Indexes(list []int) (before, after int) {
	before = slices.Index(list, self.Before)
	after = slices.Index(list, self.After)
	return
}

func parsePair(line []byte) (Pair, error) {
	before, after, found := bytes.Cut(line, []byte{'|'})
	if !found {
		return Pair{}, errors.New("Not found")
	}

	var pair Pair
	nbefore, err := strconv.ParseInt(string(before), 10, 0)
	if err != nil {
		return Pair{}, err
	}
	pair.Before = int(nbefore)
	nafter, err := strconv.ParseInt(string(after), 10, 0)
	if err != nil {
		return Pair{}, err
	}
	pair.After = int(nafter)

	return pair, nil
}

func (self *Day) D5() error {

	total := 0
	orders := []Pair{}
	lists := [][]int{}

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

		pair, err := parsePair(line)
		if err != nil {
			return err
		}
		orders = append(orders, pair)
	}
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

		list := bytes.Split(line, []byte{','})
		output := make([]int, len(list))
		for i, val := range list {
			n, err := strconv.ParseInt(string(bytes.TrimSpace(val)), 10, 0)
			if err != nil {
				return err
			}
			output[i] = int(n)
		}
		lists = append(lists, output)
	}

	toOrder := [][]int{}

	for _, list := range lists {
		good := true
		for _, pair := range orders {
			if !pair.InOrder(list) {
				good = false
				break
			}
		}
		if good {
			total += list[len(list)/2]
			// total += 1
		} else {
			toOrder = append(toOrder, list)
		}
	}

	fmt.Println(total)

	total2 := 0
	for _, list := range toOrder {
		set := PageSet{
			Pages: list,
			Rules: []Pair{},
		}
		for _, pair := range orders {
			if pair.ContainsBoth(list) {
				set.Rules = append(set.Rules, pair)
			}
		}

		set.Order()
		total2 += list[len(list)/2]
	}
	fmt.Println(total2)

	return nil
}
