package main

import "fmt"
// registers are stack pointer, program counter, data and string pool



type Vm struct {
	sp Stack[int]
	pc int
	data []string
	instructions []string
}

func (v *Vm) parseInstructions(i []string) {
	for _, input := range(i) {
		var pc int
		var instruction, t string
		_, err := fmt.Sscanf(input, "%5d %s %s", &pc, &instruction, &t)
		if err != nil {
			_, err := fmt.Sscanf(input, "%5d %s", &pc, &instruction)
			if err != nil {
				fmt.Println("ERROR parsing", input)
			}
		}
		fmt.Println(pc, instruction, t)
	}
}

func (v *Vm) run(dataSize int) {

	for {
		
	}
	
}