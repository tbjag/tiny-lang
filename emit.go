package main

import (
	"fmt"
	"strings"
)

var tokenInstruction = map[TokenType]string{
	Add:          "add",
	Subtract:     "sub",
	Multiply:     "mul",
	Divide:       "div",
	Mod:          "mod",
	Less:         "lt",
	Greater:      "gt",
	LessEqual:    "le",
	GreaterEqual: "ge",
	Equal:        "eq",
	NotEqual:     "ne",
	And:          "and",
	Or:           "or",
}

type Generator struct {
	stck         int
	variables    map[string]int
	constants    map[string]int
	instructions []string
}

func gen(n *Node) {
	g := Generator{
		stck:         0,
		variables:    make(map[string]int),
		constants:    make(map[string]int),
		instructions: make([]string, 0),
	}
	g.generate(n)
	fmt.Println("constants: ", g.constants)
	fmt.Println("variables: ", g.variables)
	for _, instruction := range g.instructions {
		fmt.Println(instruction)
	}
	fmt.Printf("%3d halt\n", g.stck)
}

func (g *Generator) generate(n *Node) {
	if n == nil {
		return
	}
	switch n.Value.Type {
	case Identifier:
		// if it isnt in map then throw error
		i := fmt.Sprintf("%3d fetch [%d]", g.stck, g.variables[n.Value.Val])
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case String, Integer:
		constant := string(n.Value.Val)
		g.constants[constant] = len(g.constants)
		i := fmt.Sprintf("%3d push %s - assigning in map %d", g.stck, constant, len(g.constants)-1)
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case Assign:
		// if variable already present then refer to old var,
		// if it isnt present create new var
		v := g.stck
		g.generate(n.Right)
		g.variables[n.Left.Value.Val] = v
		i := fmt.Sprintf("%3d store [%d]", g.stck, v)
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case Sequence:
		g.generate(n.Left)
		g.generate(n.Right)
	case Equal, NotEqual, Add, Subtract, Multiply, Divide, Mod, And, Or, Less, LessEqual, Greater, GreaterEqual:
		op := n.Value.Type
		g.generate(n.Left)
		g.generate(n.Right)
		i := fmt.Sprintf("%3d %s", g.stck, tokenInstruction[op])
		g.instructions = append(g.instructions, i)
		g.stck += 1
	case Prtc, Prti, Prts, Not, Negate:
		opCode := n.Value.Type
		g.generate(n.Left)
		i := fmt.Sprintf("%3d %s", g.stck, strings.ToLower(string(opCode)))
		g.instructions = append(g.instructions, i)
		g.stck += 1
	case If:
		g.generate(n.Left)
		jzAddress := g.stck
		jzIdx := len(g.instructions)
		g.instructions = append(g.instructions, "JZ IF UNREACHABLE")
		g.stck += 5
		// if n.Right.Value
		g.generate(n.Right)
		
		bpatch := fmt.Sprintf("%3d jz (%d) %d", jzAddress, g.stck-jzAddress-1, g.stck)
		g.instructions[jzIdx] = bpatch
	case While:
		topAddress := g.stck
		g.generate(n.Left)
		jzAddress := g.stck
		jzIdx := len(g.instructions)
		g.instructions = append(g.instructions, "JZ WHILE UNREACHABLE")
		g.stck += 5
		g.generate(n.Right)
		g.instructions = append(g.instructions, fmt.Sprintf("%3d jmp (%d) %d", g.stck, topAddress-g.stck-1, topAddress))
		g.stck += 5
		bpatch := fmt.Sprintf("%3d jz (%d) %d", jzAddress, g.stck-jzAddress-1, g.stck)
		g.instructions[jzIdx] = bpatch
	}
}
