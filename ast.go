package main

import (
	"fmt"
	"slices"
	"strings"
)

type Node struct {
	Value Token
	Left  *Node
	Right *Node
}

type Ast struct {
	tokens []Token
	tokenLen int
	root   *Node
	idx    int
	
}

func newNode(value Token) *Node {
	return &Node{
		Value: value,
		Left:  nil,
		Right: nil,
	}
}

var spacer = "  "

func printTree(ast *Node, marker string, level int) {
	if ast == nil {
		return
	}
	printTree(ast.Right, spacer+marker, level+1)
	ident := strings.Repeat(spacer, level)
	fmt.Printf("%s%s%s\n", ident, marker, ast.Value)
	printTree(ast.Left, spacer+marker, level+1)
}

/*
*       primary    =   Identifier
        | Integer
        | '(' expr ')'
        | ('+' | '-' | '!') primary
        ;
 */
func (a *Ast) primary() {
	// if a.idx < a.tokenLen && slices.Contains([])
	// ^^ contains func here to find tokentype
}
