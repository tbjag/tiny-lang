package main

import (
	"fmt"
)

// walk the tree in pre-order
// anytime you hit a generation rule emit an instruction
// keep track of stack ptr

func emit(n *Node) {
	generate(n, 0)
	fmt.Println("halt")
}

func generate(n *Node, stack int) int {
	if n == nil {
		return stack
	}
	switch n.Value.Type {
	case Identifier:
		fmt.Printf("fetch [%s]\n", n.Value.Val)
	case Integer:
		fmt.Println(stack, "push ", n.Value.Val)
		return stack + 5
	case String:
		fmt.Println("push strIndex?")
	case Assign:
		v := stack
		stack = generate(n.Right, 0)
		fmt.Printf("%d store [%d]\n",stack, v)
	case Sequence:
		generate(n.Left, stack)
		generate(n.Right, stack)
	case Equal: //binary operators
		opcode := n.Value.Type
		generate(n.Left, stack)
		generate(n.Right, stack)
		fmt.Println(opcode)
	case Prtc, Prti, Prts: // + unary
		opCode := n.Value.Type
		generate(n.Left,0)
		fmt.Println(opCode)
	case If:
		fmt.Println("in if")
	case While:
		top := stack
		// todo patch calls
		stack = generate(n.Left, stack)
		stack = generate(n.Right, stack)
		fmt.Printf("jmp to top %d", top)
	}

	return 0
}