package main

import (
	"fmt"
)

// walk the tree in pre-order
// anytime you hit a generation rule emit an instruction
// keep track of stack ptr

type Generator struct {
	stck int
}

func gen(n *Node) {
	g := Generator{stck: 0}
	g.emit(n)
}

func (g Generator) emit(n *Node) {
	g.generate(n)
	fmt.Println("halt")
}

func (g Generator) generate(n *Node) {
	if n == nil {
		return
	}
	switch n.Value.Type {
	case Identifier:
		fmt.Printf("fetch [%s]\n", n.Value.Val)
	case Integer:
		fmt.Println(g.stck, "push ", n.Value.Val)
		g.stck += 5
	case String:
		fmt.Println("push strIndex?")
	case Assign:
		v := g.stck
		g.generate(n.Right)
		fmt.Printf("%d store [%d]\n",g.stck, v)
	case Sequence:
		g.generate(n.Left)
		g.generate(n.Right)
	case Equal: //binary operators
		opcode := n.Value.Type
		g.generate(n.Left)
		g.generate(n.Right)
		fmt.Println(opcode)
	case Prtc, Prti, Prts: // + unary
		opCode := n.Value.Type
		g.generate(n.Left)
		fmt.Println(opCode)
	case If:
		fmt.Println("in if")
	case While:
		top := g.stck
		// todo patch calls
		fmt.Println("in while, ", n.Left.Value.Type)
		g.generate(n.Left)
		g.generate(n.Right)
		fmt.Printf("jmp to top %d\n", top)
	}
}