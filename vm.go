package main

import (
	"fmt"
	"strconv"
)

// registers are stack pointer, program counter, data and string pool

// create enum for cmd's
type InstructionType string

const (
	fetch InstructionType = "fetch"
	store InstructionType = "store"
	push  InstructionType = "push"
	jmp   InstructionType = "jmp"
	jz    InstructionType = "jz"
	add   InstructionType = "add"
	sub   InstructionType = "sub"
	mul   InstructionType = "mul"
	div   InstructionType = "div"
	mod   InstructionType = "mod"
	lt    InstructionType = "lt"
	gt    InstructionType = "gt"
	le    InstructionType = "le"
	ge    InstructionType = "ge"
	eq    InstructionType = "eq"
	ne    InstructionType = "ne"
	and   InstructionType = "and"
	or    InstructionType = "or"
	neg   InstructionType = "neg"
	not   InstructionType = "not"
	prtc  InstructionType = "prtc"
	prti  InstructionType = "prti"
	prts  InstructionType = "prts"
	halt  InstructionType = "halt"
)

const WordSize = 4

type Instruction struct {
	cnt  int
	cmd  InstructionType
	n    *int
	addr *int
}

type Vm struct {
	data         []int
	pool         []string
	instructions []Instruction
	addrToIdx map[int]int
}

func createVm(g *Generator) *Vm {
	vm := Vm{}
	vm.parseInstructions(g.instructions)
	vm.data = make([]int, len(g.variables))
	vm.pool = make([]string, len(g.strs))
	for i, s := range g.strs {
		// lexer keeps the surrounding quotes, so Unquote resolves \n, \", \\
		u, err := strconv.Unquote(s)
		if err != nil {
			panic("bad string literal: " + s)
		}
		vm.pool[i] = u
	}
	return &vm
}

func parseInt(s string) *int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic("error parsing str to in in push")
	}
	return &num
}

func (v *Vm) parseInstructions(i []string) {
	v.instructions = make([]Instruction, 0)
	for _, input := range i {
		var pc, n, addr int
		var instruction, t string
		var i Instruction

		_, err := fmt.Sscanf(input, "%5d %s (%d) %d", &pc, &instruction, &n, &addr)
		if err != nil {
			_, err := fmt.Sscanf(input, "%5d %s %s", &pc, &instruction, &t)
			if err != nil {
				_, err := fmt.Sscanf(input, "%5d %s", &pc, &instruction)
				if err != nil {
					panic("could not parse: " + input)
				}
				i = Instruction{cnt: pc, cmd: InstructionType(instruction)}
			} else {
				switch instruction {
				case "fetch", "store":
					n := t[1 : len(t)-1]
					i = Instruction{cnt: pc, cmd: InstructionType(instruction), n: parseInt(n)}
				case "push":
					i = Instruction{cnt: pc, cmd: InstructionType(instruction), n: parseInt(t)}
				}
			}
		} else {
			i = Instruction{cnt: pc, cmd: InstructionType(instruction), n: &n, addr: &addr}
		}
		v.instructions = append(v.instructions, i)
	}
	v.addrToIdx = make(map[int]int)
	for idx, ins := range v.instructions {
    v.addrToIdx[ins.cnt] = idx
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (v *Vm) run() {

	for idx, i := range v.instructions {
		fmt.Println(idx, i.cnt, i.cmd)
	}
	stackLen := 1000
	stack := make([]int, stackLen)
	stackPtr := 0
	pc := 0
	for {
		instr := v.instructions[pc]
		pc++
		switch instr.cmd {
		case fetch:
			stack[stackPtr] = v.data[*instr.n]
			stackPtr++
		case store:
			stackPtr--
			v.data[*instr.n] = stack[stackPtr]
		case push:
			stack[stackPtr] = *instr.n
			stackPtr++
		case add:
			stack[stackPtr-2] += stack[stackPtr-1]
			stackPtr--
		case sub:
			stack[stackPtr-2] -= stack[stackPtr-1]
			stackPtr--
		case mul:
			stack[stackPtr-2] *= stack[stackPtr-1]
			stackPtr--
		case div:
			stack[stackPtr-2] /= stack[stackPtr-1]
			stackPtr--
		case mod:
			stack[stackPtr-2] %= stack[stackPtr-1]
			stackPtr--
		case lt:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] < stack[stackPtr-1])
			stackPtr--
		case gt:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] > stack[stackPtr-1])
			stackPtr--
		case le:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] <= stack[stackPtr-1])
			stackPtr--
		case ge:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] >= stack[stackPtr-1])
			stackPtr--
		case eq:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] == stack[stackPtr-1])
			stackPtr--
		case ne:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] != stack[stackPtr-1])
			stackPtr--
		case and:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] != 0 && stack[stackPtr-1] != 0)
			stackPtr--
		case or:
			stack[stackPtr-2] = b2i(stack[stackPtr-2] != 0 || stack[stackPtr-1] != 0)
			stackPtr--
		case neg:
			stack[stackPtr-1] = -stack[stackPtr-1]
		case not:
			stack[stackPtr-1] = b2i(stack[stackPtr-1] == 0)
		case jmp:
			pc = v.addrToIdx[*instr.addr]
		case jz:
			stackPtr--
			if stack[stackPtr] == 0 {
				pc = v.addrToIdx[*instr.addr]
			}
		case prtc:
			stackPtr--
			fmt.Print(stack[stackPtr])
		case prti:
			stackPtr--
			fmt.Print(stack[stackPtr])
		case prts:
			stackPtr--
			fmt.Print(v.pool[stack[stackPtr]])

		case halt:
			return
		}
	}

}
