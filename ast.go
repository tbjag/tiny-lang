package main

import (
	"fmt"
	"maps"
	"slices"
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
	Or:    1,
	And:   2,
	Equal: 3, NotEqual: 3,
	Less: 4, LessEqual: 4, Greater: 4, GreaterEqual: 4,
	Add: 5, Subtract: 5,
	Multiply: 6, Divide: 6, Mod: 6,
	Identifier: 7, Integer: 7,
}

var precendenceKeys = slices.Collect(maps.Keys(precedenceMap))

var spacer = "  "

func printTree(n *Node) {
	if n == nil {
		fmt.Println(";")
		return
	}

	if slices.Contains([]TokenType{Identifier, Integer, String}, n.Value.Type) {
		fmt.Println(n.Value.Type, n.Value.Val)
		return
	}
	fmt.Println(n.Value.Type)
	printTree(n.Left)
	printTree(n.Right)

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
		fmt.Println(a.tokens[a.idx - 1], a.currType(), a.tokens[a.idx + 1])
		panic("failed to parse in expr")
	}

	for slices.Contains(precendenceKeys, a.currType()) && precedenceMap[a.currType()] >= p {
		saveTok := a.getToken()
		a.idx += 1
		q := precedenceMap[saveTok.Type] + 1

		prev := x
		x = newNode(saveTok)
		x.Left = prev
		x.Right = a.expr(q)
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
		t = newNode(a.getToken())
		a.idx += 1
		e := a.parenExpr()
		s := a.stmt()
		t.Left = e
		ifTok  := Token{
			Val: "",
			Type: If,
		}
		ifNode := newNode(ifTok)
		ifNode.Left = s
		t.Right = ifNode
		// todo else statement
	} else if a.currType() == Put {
		t = newNode(Token{Type: Prtc})
		a.idx += 1
		t.Left = a.parenExpr()
		a.expect(Semicolon)
	} else if a.currType() == Print {
		a.idx += 1
		a.expect(LeftParen)
		for {
			var e *Node
			if a.currType() == String {
				p := newNode(a.getToken())
				e = newNode(Token{Type: Prts})
				e.Left = p
				a.idx += 1
			} else {
				e = newNode(Token{Type: Prti})
				e.Left = a.expr(0)
			}

			n := newNode(Token{Type: Sequence})
			n.Left = t
			n.Right = e
			t = n
			if a.currType() != Comma {
				break
			}
			a.idx += 1
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
		a.idx += 1
		for a.currType() != RightBrace {
			a.checkEOF()
			s := Token{
				Val: "",
				Type: Sequence,
			}
			prev := t
			t = newNode(s)
			t.Left = prev
			t.Right = a.stmt()
		}
		a.expect(RightBrace)
	} else if a.currType() != EOF {
		panic("reached end of stmt")
	}
	return t
}
