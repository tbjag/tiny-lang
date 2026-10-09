package main

import (
	"fmt"
	"os"
)

// type AstType string
// const (
// 	Sequence AstType = "Sequence"
// 	AssignA AstType = "Assign"
// 	IdentifierA AstType = "Identifier"
// 	IntegerA AstType = "Integer"
// )
// 
// type Ast struct {
// 	Type AstType
// 	Left *Ast
// 	Right *Ast
// }

type Node struct {
	Value int
	Left *Node
	Right *Node
}

func NewNode(value int) *Node {
	return &Node{
		Value: value,
		Left: nil,
		Right: nil,
	}
}

func (n *Node) Insert(value int) {
	if value <= n.Value {
		if n.Left == nil {
			n.Left = NewNode(value)
		} else {
			n.Left.Insert(value)
		}
	} else {
		if n.Right == nil {
			n.Right = NewNode(value)
		} else {
			n.Right.Insert(value)
		}
	}
}

func InOrder(n *Node) {
	if n != nil {
		InOrder(n.Left)
		fmt.Printf("%d ", n.Value)
		InOrder(n.Right)
	}
}

func main() {
	allArgs := os.Args
	if len(allArgs) < 2 {
		fmt.Printf("provide file to lex")
		return
	}
	fileContent, err := os.ReadFile(allArgs[1])
	if err != nil {
		panic(err)
	}
	rawData := string(fileContent)

	// fmt.Println("Content: \n", rawData)

	lexRes := lex(rawData)
	fmt.Println(lexRes)
	// astRes := parse(lexRes)
	// for _, item := range astRes {
	// 	fmt.Printf("%+v\n", item)
	// }

	// fmt.Println(astRes.Type)

	fmt.Printf("\n\n\n")
	// start of AST testing!!
	root := NewNode(10)
    root.Insert(5)
    root.Insert(15)
    root.Insert(3)
    root.Insert(7)
	
    fmt.Print("In-order traversal: ")
    InOrder(root) // Output: 3 5 7 10 15
    fmt.Println()
	
}
