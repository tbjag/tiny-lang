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
	prtc  InstructionType = "prtc"
	prti  InstructionType = "prti"
	prts  InstructionType = "prts"
	halt  InstructionType = "halt"
)

type Instruction struct {
	cnt  int
	cmd  InstructionType
	n    *int
	addr *int
}

type Vm struct {
	data         []int
	pool []string
	instructions []Instruction
}

func createVm(g *Generator) *Vm {
	vm := Vm{}
	vm.parseInstructions(g.instructions)
	vm.data = make([]int, len(g.variables))
	vm.pool = make([]string, len(g.strs))
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
				if err != nil {panic("could not parse: " + input)}
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
}

func (v *Vm) run() {
	
	for idx, i := range(v.instructions) {
		fmt.Println(idx, i.cnt, i.cmd)
	}
	stack := Stack[int]{}
	pc := 0
	for {
		instr := v.instructions[v.pc]
		pc += 1
		switch instr.cmd {
			case fetch:
				stack.Push(*instr.n)
				pc += 4
			case store:
				
		}
	}

}
