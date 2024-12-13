package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type Game struct {
	A   Pos
	B   Pos
	Win Pos
}

func (self Pos) AddSelf(pos Pos) Pos {
	return self.Add(pos.X, pos.Y)
}

func (self Pos) Neg() Pos {
	return Pos{
		X: -self.X,
		Y: -self.Y,
	}
}

func (self Pos) OverZero() bool {
	return self.X >= 0 && self.Y >= 0
}
func (self Pos) IsZero() bool {
	return self.X == 0 && self.Y == 0
}

func (self Pos) IsDivisible(pos Pos) bool {
	return self.X%pos.X == 0 && self.Y%pos.Y == 0
}

func (self *Game) CalcGreedy(pos Pos) int {
	goodCost := 0
	goodState := pos
	cost := 0
	for {
		if self.Win.AddSelf(pos.AddSelf(self.B).Neg()).OverZero() {
			pos = pos.AddSelf(self.B)
			cost += 1
		} else {
			break
		}
		if self.Win.AddSelf(pos.Neg()).IsDivisible(self.B) ||
			self.Win.AddSelf(pos.Neg()).IsDivisible(self.A) {
			goodCost = cost
			goodState = pos
		}
	}

	pos = goodState
	cost = goodCost

	for {
		if self.Win.AddSelf(pos.AddSelf(self.A).Neg()).OverZero() {
			pos = pos.AddSelf(self.A)
			cost += 3
		} else {
			break
		}
	}

	if self.Win.AddSelf(pos.Neg()).IsZero() {
		return cost
	}

	return -1
}

func (self *Game) CalcFull(pos Pos, data ...int) int {
	var cost int
	if len(data) > 0 {
		cost = data[0]
	}
	var turns int
	if len(data) > 1 {
		turns = data[1]
	}
	var minCost int
	if len(data) > 2 {
		minCost = data[2]
	} else {
		minCost = self.CalcGreedy(pos)
		if minCost == -1 {
			return -1
		}
	}

	if !self.Win.AddSelf(pos.Neg()).OverZero() {
		return -1
	} else if self.Win.AddSelf(pos.Neg()).IsZero() {
		return cost
	} else if cost > minCost {
		return -1
	}

	if turns > 100 {
		return -1
	}

	aCost := self.CalcFull(pos.AddSelf(self.A), cost+3, turns+1)
	bCost := self.CalcFull(pos.AddSelf(self.B), cost+1, turns+1)

	if aCost != -1 && bCost != -1 {
		return min(aCost, bCost)
	} else if aCost != -1 {
		return aCost
	} else if bCost != -1 {
		return bCost
	}

	return minCost
}

func (self *Day) D13() error {

	games := []Game{}

	readLine := func() (string, error) {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		line = bytes.TrimSpace(line)
		return string(line), nil
	}

	var err error
	for {
		var lineA string
		for {
			lineA, err = readLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if lineA != "" {
				break
			}
		}
		lineB, err := readLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		lineC, err := readLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		parseNum := func(val string) int {
			numStr := []byte{}
			for _, b := range []byte(val) {
				if b >= '0' && b <= '9' {
					numStr = append(numStr, b)
				}
			}

			n, err := strconv.Atoi(string(numStr))
			if err != nil {
				panic(err)
			}

			return n

		}

		parsePos := func(line string, offset int) Pos {
			fields := strings.Fields(line)
			x := parseNum(fields[offset])
			y := parseNum(fields[offset+1])
			return Pos{X: x, Y: y}
		}

		games = append(games, Game{
			A:   parsePos(lineA, 2),
			B:   parsePos(lineB, 2),
			Win: parsePos(lineC, 1),
		})
	}
	total := 0

	for i, game := range games {

		mat := [][]float64{
			{float64(game.A.X), float64(game.B.X), float64(game.Win.X)},
			{float64(game.A.Y), float64(game.B.Y), float64(game.Win.Y)},
		}

		aDiv := mat[0][0] / mat[1][0]

		mat[1] = []float64{
			mat[1][0] * aDiv,
			mat[1][1] * aDiv,
			mat[1][2] * aDiv,
		}

		mat[1] = []float64{
			mat[0][0] - mat[1][0],
			mat[0][1] - mat[1][1],
			mat[0][2] - mat[1][2],
		}

		bDiv := mat[1][1] / mat[0][1]
		mat[0] = []float64{
			mat[0][0] * bDiv,
			mat[0][1] * bDiv,
			mat[0][2] * bDiv,
		}
		mat[0] = []float64{
			mat[1][0] - mat[0][0],
			mat[1][1] - mat[0][1],
			mat[1][2] - mat[0][2],
		}

		a := mat[0][2] / mat[0][0]
		b := mat[1][2] / mat[1][1]

		isInt := func(val float64) bool {
			val = math.Round(val*100000) / 100000
			return math.Abs(val-float64(int(val))) < 0.001
		}

		if isInt(a) &&
			isInt(b) {
			ai := int(math.Round(a*100000) / 100000)
			bi := int(math.Round(b*100000) / 100000)
			if ai < 0 || bi < 0 {
				Debugf("Invalid Scores %d: %d,%d\n", i, ai, bi)
				continue
			}
			total += ai*3 + bi
		} else {
			Debugf("Game %d: Unsolveable %f,%f(%f)\n", i, a, b, a+b)
		}

	}

	total2 := 0
	for i, game := range games {

		mat := [][]float64{
			{float64(game.A.X), float64(game.B.X), float64(game.Win.X + 10000000000000)},
			{float64(game.A.Y), float64(game.B.Y), float64(game.Win.Y + 10000000000000)},
		}

		aDiv := mat[0][0] / mat[1][0]

		mat[1] = []float64{
			mat[1][0] * aDiv,
			mat[1][1] * aDiv,
			mat[1][2] * aDiv,
		}

		mat[1] = []float64{
			mat[0][0] - mat[1][0],
			mat[0][1] - mat[1][1],
			mat[0][2] - mat[1][2],
		}

		bDiv := mat[1][1] / mat[0][1]
		mat[0] = []float64{
			mat[0][0] * bDiv,
			mat[0][1] * bDiv,
			mat[0][2] * bDiv,
		}
		mat[0] = []float64{
			mat[1][0] - mat[0][0],
			mat[1][1] - mat[0][1],
			mat[1][2] - mat[0][2],
		}

		a := mat[0][2] / mat[0][0]
		b := mat[1][2] / mat[1][1]

		isInt := func(val float64) bool {
			val = math.Round(val*1000) / 1000
			return math.Abs(val-float64(int(val))) < 0.001
		}

		if isInt(a) &&
			isInt(b) {
			ai := int(math.Round(a*1000) / 1000)
			bi := int(math.Round(b*1000) / 1000)
			if ai < 0 || bi < 0 {
				Debugf("Invalid Scores %d: %d,%d\n", i, ai, bi)
				continue
			}
			total2 += ai*3 + bi
		} else {
			Debugf("Game %d: Unsolveable %f,%f(%f)\n", i, a, b, a+b)
		}

	}

	fmt.Println(total)
	fmt.Println(total2)

	return nil
}
