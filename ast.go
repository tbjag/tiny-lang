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
		panic("reached EOF")
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

func (a *Ast) expect(t TokenType) {
	if a.currType() != t {
		panic(fmt.Sprintf("Expected %s but got %s", t, a.currType()))
	}
	a.idx += 1
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
	a.expect(LeftParen)
	x := a.expr(0)
	a.expect(RightParen)
	return x
}

func (a *Ast) stmt() *Node {
	var t *Node
	if a.currType() == If {
		curr := a.getToken()
		e := a.parenExpr()
        s := a.stmt()
        internal := newNode(curr)
        a.idx += 1
        internal.Left = s
        if a.currType() == Else {
        	t.Right = a.stmt()
        }
        t = newNode(curr)
        t.Left = e
        t.Right = internal
	} else if a.currType() == Put {
        t = newNode(a.getToken())
        a.idx += 1
        t.Left = a.parenExpr()
        a.expect(Semicolon)
	} else if a.currType() == Print {
		printToken := a.getToken()
		a.idx += 1
		a.expect(LeftParen)
		for {
			var e *Node
			if a.currType() == String {
				p := newNode(a.getToken())
				e = newNode(printToken)
				e.Left = p
				a.idx += 1
			} else {
				e = a.expr(0)
			}
			seq := Token{
				Type: Semicolon,
				Val: "",
			}
			t = newNode(seq) //??
			t.Left = t //??
			t.Right = e
			if a.currType() != Comma {
				break
			}
		}
		a.expect(RightParen)
		a.expect(Semicolon)
	} else if a.currType() == Semicolon {
		a.idx += 1
	} else if a.currType() == Identifier {
		v := newNode(a.getToken())
		a.idx += 1
		op := a.getToken()
		a.expect(Assign)
		t = newNode(op)
		t.Left = v
		t.Right = a.expr(0)
		a.expect(Semicolon)
	} else if a.currType() == While {
		t = newNode(a.getToken())
		a.idx += 1
		e := a.parenExpr()
		t.Left = e
		t.Right = a.stmt()
	} else if a.currType() == LeftBrace {
		panic("todo")
	} else if a.currType() != EOF {
		panic("reached end of stmt")
	}
	return t
}
