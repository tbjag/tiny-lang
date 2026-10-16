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

func newNode(value Token) *Node {
	return &Node{
		Value: value,
		Left:  nil,
		Right: nil,
	}
}

type Ast struct {
	tokens   []Token
	tokenLen int
	root     *Node
	idx      int
}

var precedenceMap = map[TokenType]int{
	Subtract: 1,
	Add:      2,
	Not:      3,
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

func (a *Ast) checkEOF() {
	if a.tokens[a.idx].Type == EOF {
		panic("idx greater than token length")
	}
}

func (a *Ast) currType() TokenType {
	a.checkEOF()
	return a.tokens[a.idx].Type
}

func (a *Ast) getToken() Token {
	a.checkEOF()
	return a.tokens[a.idx]
}

func (a *Ast) expr(p int) *Node {
	var x *Node
	if a.currType() == LeftParen {
		x = a.parenExpr()
	} else if slices.Contains([]TokenType{Subtract, Add, Not}, a.currType()) {
		op := a.getToken()
		a.idx += 1
		y := a.expr(precedenceMap[a.currType()])
		if op.Type == Add {
			x = y
		} else {
			x = newNode(op)
			x.Left = y
		}
	} else if slices.Contains([]TokenType{Identifier, Integer}, a.currType()) {
		x = newNode(a.getToken())
		a.idx += 1
	} else {
		panic("failed to parse in expr")
	}

	for slices.Contains([]TokenType{Add, Subtract, Not}, a.currType()) && precedenceMap[a.currType()] >= p {
		saveTok := a.getToken()
		a.idx += 1
		q := precedenceMap[a.currType()]
		// if saveTok != waht?
		x = newNode(saveTok)
		x.Left = a.expr(q)
	}

	return x
}

func (a *Ast) parenExpr() *Node {
	if a.currType() != LeftParen {
		panic("expected leftParen")
	}
	x := a.expr(0)
	if a.currType() != RightParen {
		panic("expected RightParen")
	}
	return x
}

func (a *Ast) stmt() *Node {
	var t *Node
	if a.currType() == If {
		panic("todo")
	} else if a.currType() == Put {
		panic("todo")
	} else if a.currType() == Print {
		panic("todo")
	} else if a.currType() == Semicolon {
		a.idx += 1
	} else if a.currType() == Identifier {
		panic("todo")
	} else if a.currType() == While {
		panic("todo")
	} else if a.currType() == LeftBrace {
		panic("todo")
	} else if a.currType() != EOF {
		panic("reached end of stmt")
	}
	return t
}
