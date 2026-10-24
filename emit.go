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
	constants    int
	literals     map[string]int
	strs         []string
	instructions []string
}

type JAddress string

const (
	JMP = "jmp"
	JZ  = "jz"
)

type Hole struct {
	jAddress JAddress
	idx      int
	address  int
}

func gen(n *Node) string {
	g := Generator{
		stck:         0,
		variables:    make(map[string]int),
		constants:    0,
		literals:     make(map[string]int),
		instructions: make([]string, 0),
	}
	g.generate(n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "Datasize: %d Strings: %d\n", len(g.variables), len(g.strs))
	for _, l := range g.strs {
		sb.WriteString(l + "\n")
	}
	for _, instruction := range g.instructions {
		sb.WriteString(instruction + "\n")
	}
	fmt.Fprintf(&sb, "%5d halt\n", g.stck)
	return sb.String()
}

func (g *Generator) patchTo(h Hole, target int) {
	g.instructions[h.idx] = fmt.Sprintf("%5d %-6s (%d) %d", h.address, h.jAddress, target-h.address-1, target)
}

func (g *Generator) patch(h Hole) {
	g.patchTo(h, g.stck)
}

func (g *Generator) hole(j JAddress) Hole {
	h := Hole{
		jAddress: j,
		idx:      len(g.instructions),
		address:  g.stck,
	}
	g.instructions = append(g.instructions, "INSTR HOLE UNREACHABLE")
	g.stck += 5
	return h
}

func (g *Generator) generate(n *Node) {
	if n == nil {
		return
	}
	switch n.Value.Type {
	case Identifier:
		// if it isnt in map then throw error
		i := fmt.Sprintf("%5d fetch [%d]", g.stck, g.variables[n.Value.Val])
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case String:
		idx, ok := g.literals[n.Value.Val]
		if !ok {
			idx = len(g.strs)
			g.literals[n.Value.Val] = idx
			g.strs = append(g.strs, n.Value.Val)
		}
		i := fmt.Sprintf("%5d push %2d", g.stck, idx)
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case Integer:
		constant := n.Value.Val
		g.constants += 1
		i := fmt.Sprintf("%5d push  %s", g.stck, constant)
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case Character:
		i := fmt.Sprintf("%5d push  %d", g.stck, charCode(n.Value.Val))
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case Assign:
		v, ok := g.variables[n.Left.Value.Val]
		if !ok {
			v = len(g.variables)
			g.variables[n.Left.Value.Val] = v
		}
		g.generate(n.Right)
		i := fmt.Sprintf("%5d store [%d]", g.stck, v)
		g.instructions = append(g.instructions, i)
		g.stck += 5
	case Sequence:
		g.generate(n.Left)
		g.generate(n.Right)
	case Equal, NotEqual, Add, Subtract, Multiply, Divide, Mod, And, Or, Less, LessEqual, Greater, GreaterEqual:
		op := tokenInstruction[n.Value.Type]
		if n.Right == nil { // unary minus
			op = "neg"
		}
		g.generate(n.Left)
		g.generate(n.Right)
		i := fmt.Sprintf("%5d %s", g.stck, op)
		g.instructions = append(g.instructions, i)
		g.stck += 1
	case Prtc, Prti, Prts, Not, Negate:
		opCode := n.Value.Type
		g.generate(n.Left)
		i := fmt.Sprintf("%5d %s", g.stck, strings.ToLower(string(opCode)))
		g.instructions = append(g.instructions, i)
		g.stck += 1
	case If:
		g.generate(n.Left)
		jzHole := g.hole(JZ)
		g.generate(n.Right.Left) // then

		if n.Right.Right != nil { // else
			jmpHole := g.hole(JMP)
			g.patch(jzHole)
			g.generate(n.Right.Right)
			g.patch(jmpHole)
		} else {
			g.patch(jzHole)
		}
	case While:
		start := g.stck
		g.generate(n.Left)
		jzHole := g.hole(JZ)
		g.generate(n.Right)
		jmpHole := g.hole(JMP)
		g.patchTo(jmpHole, start)
		g.patch(jzHole)
	}
}