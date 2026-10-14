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

func (a *Ast) checkEOF() {
	if a.tokens[a.idx].Type == EOF {
		panic("idx greater than token length")
	}
}

/*
*       primary    =   Identifier
        | Integer
        | '(' expr ')'
        | ('+' | '-' | '!') primary
        ;
 */
func (a *Ast) primary() *Node { //TODO need support primary
	a.checkEOF()
	currentToken := a.tokens[a.idx]
	a.idx += 1
	if slices.Contains([]TokenType{Integer, Identifier}, currentToken.Type) {
		return newNode(currentToken)
	}
	if currentToken.Type == LeftParen {
		expr := a.expression()
	 	a.checkEOF()
		curr := a.tokens[a.idx]
		if curr.Type == RightParen {
			a.idx += 1
			return expr
		} else {
			panic("expected right paren in terminal")
		}
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
	a.checkEOF()
	for slices.Contains([]TokenType{Multiply, Divide, Mod}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		n := newNode(op)
		n.Left = res
		n.Right = a.primary()
		res = n
	}
	return res
}

/*
 * addition_expr = multiplication_expr {('+' | '-') multiplication_expr} ;
 */
func (a *Ast) additionExpression() *Node {
	res := a.multiplicationExpression()
	a.checkEOF()
	for slices.Contains([]TokenType{Add, Subtract}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		n := newNode(op)
		n.Left = res
		n.Right = a.multiplicationExpression()
		res = n
	}
	return res
}

/*
 * relational_expr = addition_expr [('<' | '<=' | '>' | '>=') addition_expr] ;
 */
func (a *Ast) relationalExpression() *Node {
	res := a.additionExpression()
	a.checkEOF()
	for slices.Contains([]TokenType{Less, LessEqual, Greater, GreaterEqual}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		n := newNode(op)
		n.Left = res
		n.Right = a.additionExpression()
		res = n
	}
	return res
}

/*
 * equality_expr = relational_expr [('==' | '!=') relational_expr] ;
 */
func (a *Ast) equalityExpression() *Node {
	res := a.relationalExpression()
	a.checkEOF()
	for slices.Contains([]TokenType{Equal, NotEqual}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		n := newNode(op)
		n.Left = res
		n.Right = a.relationalExpression()
		res = n
	}
	return res
}

/*
 * and_expr = equality_expr {'&&' equality_expr} ;
 */
func (a *Ast) andExpression() *Node {
	res := a.equalityExpression()
	a.checkEOF()
	for slices.Contains([]TokenType{And}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		n := newNode(op)
		n.Left = res
		n.Right = a.equalityExpression()
		res = n
	}
	return res
}

/*
 * expr = and_expr {'||' and_expr} ;
 */
func (a *Ast) expression() *Node {
	res := a.andExpression()
	a.checkEOF()
	for slices.Contains([]TokenType{Or}, a.tokens[a.idx].Type) {
		op := a.tokens[a.idx]
		a.idx += 1
		n := newNode(op)
		n.Left = res
		n.Right = a.andExpression()
		res = n
	}
	return res
}

/*
 * prt_list = (string | expr) {',' (String | expr)} ;
 */
// func (a *Ast) prtList() *Node {
// 	res := a.expression()
// 	a.checkEOF()
// 	fmt.Println("GOT")
// 	for slices.Contains([]TokenType{Comma}, a.tokens[a.idx].Type) {
// 		op := a.tokens[a.idx]
// 		a.idx += 1
// 		n := newNode(op)
// 		n.Left = res
// 		n.Right = a.expression()
// 		res = n
// 	}
// 	return res
// }