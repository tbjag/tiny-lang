package main

import (
	"fmt"
	"strings"
)

// walk the tree in pre-order
// anytime you hit a generation rule emit an instruction
// keep track of stack ptr
// create string assignments on the fly
// create a map of token -> instruction
var tokenInstruction = map[TokenType]string{
	Less: "lt",
	Add: "add",
}

type Generator struct {
	stck int
	data map[string]int
}

func gen(n *Node) {
	g := Generator{stck: 0, data: make(map[string]int)}
	g.emit(n)
}

func (g *Generator) emit(n *Node) {
	g.generate(n)
	fmt.Println("halt")
}

func (g *Generator) generate(n *Node) {
	if n == nil {
		return
	}
	switch n.Value.Type {
	case Identifier:
		// if it isnt in map then throw error
		fmt.Printf("%d fetch [%d]\n", g.stck, g.data[n.Value.Val])
		g.stck += 5
	case Integer:
		fmt.Println(g.stck, "push ", n.Value.Val)
		g.stck += 5
	case String:
		fmt.Printf("%d push %s\n", g.stck, n.Value.Val)
		g.stck += 5
	case Assign:
		v := g.stck
		g.generate(n.Right)
		g.data[n.Left.Value.Val] = v
		fmt.Printf("%d store [%d]\n", g.stck, v)
		g.stck += 5
	case Sequence:
		g.generate(n.Left)
		g.generate(n.Right)
	case Equal, NotEqual, Add, Subtract, Multiply, Divide, Mod, And, Or, Less, LessEqual, Greater, GreaterEqual:
		op := n.Value.Type
		g.generate(n.Left)
		g.generate(n.Right)
		fmt.Printf("%d %s\n", g.stck, tokenInstruction[op])
		g.stck += 1
	case Prtc, Prti, Prts, Not, Negate:
		opCode := n.Value.Type
		g.generate(n.Left)
		fmt.Printf("%d %s\n", g.stck, strings.ToLower(string(opCode)))
		g.stck += 1
	case If:
		fmt.Println("in if")
	case While:
		topAddress := g.stck
		// todo patch calls
		g.generate(n.Left)
		fmt.Printf("%d jz (xxx) xxx\n", g.stck)
		g.stck += 5
		g.generate(n.Right)
		fmt.Printf("%d jmp (xxx) %d\n", g.stck, topAddress)
	}
}
