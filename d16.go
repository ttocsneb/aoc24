package main

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"time"
)

const (
	CStart rune = 'S'
	CEnd   rune = 'E'
	CWall  rune = '#'
	CEmpty rune = '.'
	CUp    rune = '^'
	CDown  rune = 'v'
	CLeft  rune = '<'
	CRight rune = '>'
)

type PositionState struct {
	Pos Pos
	Dir rune
}

type MapState struct {
	Pos      Pos
	Dir      rune
	Previous []*MapState
	Next     []*MapState
	score    int
}

func PosesContains(poses []Pos, target Pos) (int, bool) {
	return slices.BinarySearchFunc(poses, target, ComparePos)
}

func (self *MapState) State() PositionState {
	return PositionState{
		Pos: self.Pos,
		Dir: self.Dir,
	}
}

func (self *MapState) NextStates(walls []Pos) []*MapState {
	states := make([]*MapState, 0, 3)
	var ccw rune
	var cw rune
	mov := Pos{X: 0, Y: 0}
	switch self.Dir {
	case CUp:
		mov.Y = -1
		ccw = CLeft
		cw = CRight
	case CDown:
		mov.Y = 1
		ccw = CRight
		cw = CLeft
	case CLeft:
		mov.X = -1
		ccw = CUp
		cw = CDown
	case CRight:
		mov.X = 1
		ccw = CDown
		cw = CUp
	}
	movPos := self.Pos.AddSelf(mov)
	if _, found := PosesContains(walls, movPos); !found {
		states = append(states, &MapState{
			Pos:      movPos,
			Dir:      self.Dir,
			Previous: []*MapState{self},
			score:    1,
			Next:     []*MapState{},
		})
	}

	states = append(states, &MapState{
		Pos:      self.Pos,
		Dir:      ccw,
		Previous: []*MapState{self},
		score:    1000,
		Next:     []*MapState{},
	}, &MapState{
		Pos:      self.Pos,
		Dir:      cw,
		Previous: []*MapState{self},
		score:    1000,
		Next:     []*MapState{},
	})

	for _, next := range states {
		self.Next = append(self.Next, next)
	}

	return states
}

func ComparePosState(lhs, rhs PositionState) int {
	if cmp := ComparePos(lhs.Pos, rhs.Pos); cmp != 0 {
		return cmp
	}
	if lhs.Dir < rhs.Dir {
		return 1
	}
	if lhs.Dir > rhs.Dir {
		return -1
	}
	return 0
}

func IndexPosState(states []*MapState, wanted PositionState) int {
	return slices.IndexFunc(states, func(state *MapState) bool {
		return ComparePosState(state.State(), wanted) == 0
	})
}

func (self *MapState) SwapNexts(newParent *MapState) {
	newScore := newParent.Score()
	for _, next := range self.Next {
		if i := slices.Index(next.Previous, self); i != -1 {
			next.Previous[i] = newParent
		}

		i := 0
		for i < len(next.Previous) {
			sibling := next.Previous[i]
			if sibling == newParent {
				i++
				continue
			}
			if sibling.Score() > newScore {
				next.Previous = slices.Delete(next.Previous, i, i+1)
			} else {
				i++
			}
		}
	}
	self.Next = []*MapState{}
}

func (self *MapState) Score() int {
	score := self.score
	if len(self.Previous) > 0 {
		score += self.Previous[0].Score()
	}
	return score
}

func (self *Day) D16() error {

	warehouse := [][]rune{}

	readLine := func() (string, error) {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		line = bytes.TrimSpace(line)
		return string(line), nil
	}
	firstLine, err := readLine()
	if err != nil {
		return err
	}
	warehouse = append(warehouse, []rune(firstLine))

	// Read warehouse
	for {
		line, err := readLine()
		if err != nil {
			return err
		}
		warehouse = append(warehouse, []rune(line))

		if line == firstLine {
			break
		}
	}

	walls := []Pos{}
	var start Pos
	var end Pos

	for y, row := range warehouse {
		for x, c := range row {
			p := Pos{X: x, Y: y}
			if c == CWall {
				if i, found := PosesContains(walls, p); !found {
					walls = slices.Insert(walls, i, p)
				}
			} else if c == CStart {
				start = p
			} else if c == CEnd {
				end = p
			}
		}
	}
	if start.X == 0 && start.Y == 0 || end.X == 0 || end.Y == 0 {
		return errors.New("Could not find one of start or end")
	}

	Debug("Start: ", start)
	Debug("End: ", end)

	nextStates := []*MapState{{
		Pos:      start,
		Dir:      CRight,
		Previous: []*MapState{},
		Next:     []*MapState{},
		score:    0,
	}}

	seenStates := map[PositionState]*MapState{}

	last := time.Now()
	lastCount := 0
	hadBetter := false
	runPrediction := func() {
		fmt.Print("0")
		done := false
		hadBetter = false
		seen := 0
		for len(nextStates) > 0 {
			next := nextStates[0]
			nextStates = slices.Delete(nextStates, 0, 1)
			seen += 1
			if time.Now().Sub(last) > time.Millisecond*50 {
				if done {
					Debug()
				}
				fmt.Printf("\033[2K\r%d seen, %d left, %d/s", seen, len(nextStates), int(float64(len(nextStates)-lastCount)/0.050))
				done = false
				lastCount = len(nextStates)
				last = time.Now()
			}

			for _, possible := range next.NextStates(walls) {
				s := possible.State()
				if existing, found := seenStates[s]; found {
					lhs := existing.Score()
					rhs := possible.Score()
					if lhs > rhs {
						existing.SwapNexts(possible)
						seenStates[s] = possible
						hadBetter = true
					} else if lhs == rhs {
						if i := IndexPosState(existing.Previous, next.State()); i != -1 {
							existing.Previous[i] = next
						} else {
							existing.Previous = append(existing.Previous, next)
						}
					}
				} else {
					nextStates = append(nextStates, possible)
					seenStates[s] = possible
				}
			}
		}
		fmt.Println()
	}

	score := 99999999999
	bestPaths := []*MapState{}
	bestPrediction := func() {
		if state, found := seenStates[PositionState{Pos: end, Dir: CUp}]; found {
			s := state.Score()
			if score > s {
				bestPaths = []*MapState{}
				score = s
			}
			if score == s {
				if !slices.Contains(bestPaths, state) {
					bestPaths = append(bestPaths, state)
				}
			}
		}
		if state, found := seenStates[PositionState{Pos: end, Dir: CDown}]; found {
			s := state.Score()
			if score > s {
				bestPaths = []*MapState{}
				score = s
			}
			if score == s {
				if !slices.Contains(bestPaths, state) {
					bestPaths = append(bestPaths, state)
				}
			}
		}
		if state, found := seenStates[PositionState{Pos: end, Dir: CLeft}]; found {
			s := state.Score()
			if score > s {
				bestPaths = []*MapState{}
				score = s
			}
			if score == s {
				if !slices.Contains(bestPaths, state) {
					bestPaths = append(bestPaths, state)
				}
			}
		}
		if state, found := seenStates[PositionState{Pos: end, Dir: CRight}]; found {
			s := state.Score()
			if score > s {
				bestPaths = []*MapState{}
				score = s
			}
			if score == s {
				if !slices.Contains(bestPaths, state) {
					bestPaths = append(bestPaths, state)
				}
			}
		}

		for _, state := range seenStates {
			nextStates = append(nextStates, state)
		}

	}
	lastScore := 0
	count := 0
	for {
		runPrediction()
		bestPrediction()
		if !hadBetter {
			Debug("No more better")
			break
		}
		if lastScore == score {
			count += 1
			if count > 5 {
				break
			}
		} else {
			count = 0
		}
		lastScore = score
		fmt.Printf("Best score so far: %d\n", score)
	}

	fmt.Println(score)

	path := []Pos{}
	nPaths := 1
	for _, best := range bestPaths {
		toCheck := []*MapState{best}
		for len(toCheck) > 0 {
			state := toCheck[0]
			toCheck = toCheck[1:]
			if !slices.Contains(path, state.Pos) {
				path = append(path, state.Pos)
				warehouse[state.Pos.Y][state.Pos.X] = 'O'
			}

			if len(state.Previous) != 0 {
				nPaths += len(state.Previous) - 1
			}
			for _, prev := range state.Previous {
				toCheck = append(toCheck, prev)
			}
		}
	}

	for _, row := range warehouse {
		Debug(string(row))
	}
	Debugf("There are %d branches\n", nPaths)

	fmt.Println(len(path))

	return nil
}
