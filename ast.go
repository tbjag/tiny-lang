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

func newNodeTok(t Token) *Node {
	return &Node{
		Value: t,
		Left:  nil,
		Right: nil,
	}
}

func newNodeTyp(t TokenType) *Node {
	return &Node{
		Value: Token{Type: t},
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

func createAst(t []Token) *Node {
	ast := Ast{
		tokens:   t,
		tokenLen: len(t),
		root:     nil,
		idx:      0,
	}

	for ast.tokens[ast.idx].Type != EOF {
		stmt := ast.stmt()
		n := newNodeTyp(Sequence)
		n.Left = ast.root
		n.Right = stmt
		ast.root = n
	}

	return ast.root
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

func printTreePreOrder(n *Node) {
	if n == nil {
		fmt.Println(";")
		return
	}
	if slices.Contains([]TokenType{Identifier, Integer, String}, n.Value.Type) {
		fmt.Println(n.Value.Type, n.Value.Val)
		return
	}
	fmt.Println(n.Value.Type)
	printTreePreOrder(n.Left)
	printTreePreOrder(n.Right)
}

func printTreeNested(n *Node) string {
	if n == nil {
		return "nil"
	}
	switch n.Value.Type {
	case Identifier, Integer, String:
		return n.Value.Val
	}
	name := string(n.Value.Type)
	if n.Left == nil && n.Right == nil {
		return name
	}
	return fmt.Sprintf("%s(%s, %s)", name, printTreeNested(n.Left), printTreeNested(n.Right))
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

func (a *Ast) accept(t TokenType) bool {
	if a.currType() == t {
		a.idx += 1
		return true
	}
	return false
}

func (a *Ast) acceptReturn(t []TokenType) (bool, *Token) {
	if slices.Contains(t, a.currType()) {
		a.idx += 1
		return true, &(a.tokens[a.idx-1])
	}
	return false, nil
}

func (a *Ast) isBinary(p int) bool {
	prec, ok := precedenceMap[a.tokens[a.idx].Type]
	return ok && prec >= p
}

var ExprOperands = []TokenType{Subtract, Add, Not}
var AssignOperands = []TokenType{Identifier, Integer}

func (a *Ast) expr(p int) *Node {
	var x *Node
	if a.currType() == LeftParen {
		x = a.parenExpr()
	} else if acc, op := a.acceptReturn(ExprOperands); acc {
		y := a.expr(precedenceMap[a.currType()])
		if op.Type == Add {
			x = y
		} else {
			x = newNodeTok(*op)
			x.Left = y
		}
	} else if acc, op := a.acceptReturn(AssignOperands); acc {
		x = newNodeTok(*op)
	} else {
		fmt.Println(a.tokens[a.idx-1], a.currType(), a.tokens[a.idx+1]) // todo better print formatting for failed parse
		panic("failed to parse in expr")
	}

	for a.isBinary(p) {
		saveTok := a.getToken()
		a.idx += 1
		q := precedenceMap[saveTok.Type] + 1

		prev := x
		x = newNodeTok(saveTok)
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
	if acc := a.accept(If); acc {
		t = newNodeTyp(If)
		e := a.parenExpr()
		s := a.stmt()
		t.Left = e
		i := newNodeTyp(If)
		i.Left = s
		if a.currType() == Else {
			a.idx += 1
			i.Right = a.stmt()
		}
		t.Right = i
	} else if acc := a.accept(Put); acc {
		t = newNodeTyp(Prtc)
		t.Left = a.parenExpr()
		a.expect(Semicolon)
	} else if acc := a.accept(Print); acc {
		a.expect(LeftParen)
		for {
			var e *Node
			if acc, op := a.acceptReturn([]TokenType{String}); acc {
				p := newNodeTok(*op)
				e = newNodeTyp(Prts)
				e.Left = p
			} else {
				e = newNodeTyp(Prti)
				e.Left = a.expr(0)
			}
			n := newNodeTyp(Sequence)
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
	} else if acc := a.accept(Semicolon); acc {
		// expecting semicolon here
	} else if a.currType() == Identifier {
		v := newNodeTok(a.getToken())
		a.idx += 1
		t = newNodeTok(a.getToken())
		a.expect(Assign)
		t.Left = v
		t.Right = a.expr(0)
		a.expect(Semicolon)
	} else if acc := a.accept(While); acc {
		t = newNodeTyp(While)
		e := a.parenExpr()
		t.Left = e
		t.Right = a.stmt()
	} else if acc := a.accept(LeftBrace); acc {
		for a.currType() != RightBrace {
			a.checkEOF()
			prev := t
			t = newNodeTyp(Sequence)
			t.Left = prev
			t.Right = a.stmt()
		}
		a.expect(RightBrace)
	} else if a.currType() != EOF {
		panic("reached end of stmt")
	}
	return t
}
