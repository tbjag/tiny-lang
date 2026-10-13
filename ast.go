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
func (a *Ast) primary() *Node {
	if a.idx >= a.tokenLen {
		panic("idx greater than token length")
	}
	currentToken := a.tokens[a.idx]
	a.idx += 1
	if slices.Contains([]TokenType{Integer, Identifier}, currentToken.Type) {
		return newNode(currentToken)
	}
	if currentToken.Type == LeftParen {
		// call expr here
		panic("todo")
	}
	// if slices.Contains([]TokenType{Integer, Identifier}, currentToken.Type) {
	// 	return newNode(currentToken) // this is +, -, !
	// }
	panic(fmt.Sprintf("unexpected token in primary %s", currentToken.Type))
}

/*
 * multiplication_expr =   primary {('*' | '/' | '%') primary } ;
 */
func (a *Ast) multiplicationExpression() *Node {
	res := a.primary()
	if a.idx >= a.tokenLen {
		panic("idx greater than token length")
	}
	
	for slices.Contains([]TokenType{Multiply, Divide, Mod}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		if op.Type == Multiply {
			n := newNode(op)
			n.Left = res
			n.Right = a.primary()
		} else {
			panic("not implemented")
		}
	}
	return res
}

/*
 * addition_expr = multiplication_expr {('+' | '-') multiplication_expr} ;
 */
func (a *Ast) additionExpression() *Node {
	res := a.multiplicationExpression()
	if a.idx >= a.tokenLen {
		panic("idx greater than token length")
	}
}