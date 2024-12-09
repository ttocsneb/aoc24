package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Equation struct {
	Result int
	Inputs []int
}

type Op int

const (
	Add Op = iota
	Mul
	Concat
)

func (self Op) String() string {
	if self == Add {
		return "+"
	}
	if self == Mul {
		return "*"
	}
	return "?"
}

func (self *Equation) CanSolve() bool {
	n := len(self.Inputs) - 1
	possible := 1 << n

	for i := 0; i < possible; i++ {
		test := strings.Builder{}
		result := self.Inputs[0]
		if isDebug {
			test.WriteString(fmt.Sprint(result))
		}
		for opI := 0; opI < n; opI++ {
			if (1<<opI)&i == 0 {
				if isDebug {
					test.WriteString(fmt.Sprintf(" + %d", self.Inputs[opI+1]))
				}
				result = result + self.Inputs[opI+1]
			} else {
				if isDebug {
					test.WriteString(fmt.Sprintf(" * %d", self.Inputs[opI+1]))
				}
				result = result * self.Inputs[opI+1]
			}
		}
		if isDebug {
			test.WriteString(fmt.Sprintf(" = %d (%d)", result, self.Result))
		}
		if result == self.Result {
			Debug(test.String())
			return true
		}
	}
	return false
}

func incOps(ops []Op, i int) bool {
	if i >= len(ops) {
		return true
	}
	ops[i] += 1
	if ops[i] > Concat {
		ops[i] = 0
		return incOps(ops, i+1)
	}
	return false
}

func (self *Equation) CanSolveConcat() bool {
	n := len(self.Inputs) - 1
	ops := make([]Op, n)

	for {
		test := strings.Builder{}
		result := self.Inputs[0]
		if isDebug {
			test.WriteString(fmt.Sprint(result))
		}
		for i := 0; i < n; i++ {
			if ops[i] == Add {
				if isDebug {
					test.WriteString(fmt.Sprintf(" + %d", self.Inputs[i+1]))
				}
				result = result + self.Inputs[i+1]
			} else if ops[i] == Mul {
				if isDebug {
					test.WriteString(fmt.Sprintf(" * %d", self.Inputs[i+1]))
				}

				result = result * self.Inputs[i+1]
			} else if ops[i] == Concat {
				if isDebug {
					test.WriteString(fmt.Sprintf(" || %d", self.Inputs[i+1]))
				}
				v, _ := strconv.ParseInt(fmt.Sprintf("%d%d", result, self.Inputs[i+1]), 10, 0)
				result = int(v)
			}
		}
		if isDebug {
			test.WriteString(fmt.Sprintf(" = %d (%d)", result, self.Result))
		}
		if result == self.Result {
			Debug(test.String())
			return true
		}

		overflow := incOps(ops, 0)
		if overflow {
			break
		}
	}
	return false
}

func (self *Day) D7() error {
	count := 0
	p2count := 0

	for {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line = bytes.TrimSpace(line)

		res, params, _ := strings.Cut(string(line), ":")
		eq := Equation{
			Inputs: []int{},
		}
		n, err := strconv.ParseInt(res, 10, 0)
		if err != nil {
			return err
		}
		eq.Result = int(n)

		for _, field := range strings.Fields(params) {
			n, err = strconv.ParseInt(field, 10, 0)
			if err != nil {
				return err
			}
			eq.Inputs = append(eq.Inputs, int(n))
		}
		if eq.CanSolve() {
			count += eq.Result
		}
		if eq.CanSolveConcat() {
			p2count += eq.Result
		}
	}

	fmt.Println(count)
	fmt.Println(p2count)

	return nil
}
