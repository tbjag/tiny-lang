package main

type TokenType string

const (
	// Special
	EOF TokenType = "EOF"

	// Keywords
	Print TokenType = "Print"
	Put   TokenType = "Put"
	While TokenType = "While"
	If    TokenType = "If"
	Else  TokenType = "Else"

	// Literals & Identifiers
	Identifier TokenType = "Identifier"
	Integer    TokenType = "Integer"
	String     TokenType = "String"
	Character  TokenType = "Character"

	// Operators
	Assign       TokenType = "Assign"
	Add          TokenType = "Add"
	Subtract     TokenType = "Subtract"
	Multiply     TokenType = "Multiply"
	Divide       TokenType = "Divide"
	Mod          TokenType = "Mod"
	Negate       TokenType = "Negate"
	Not          TokenType = "Not"
	Less         TokenType = "Less"
	LessEqual    TokenType = "LessEqual"
	Greater      TokenType = "Greater"
	GreaterEqual TokenType = "GreaterEqual"
	Equal        TokenType = "Equal"
	NotEqual     TokenType = "NotEqual"
	And          TokenType = "And"
	Or           TokenType = "Or"

	// Delimiters
	LeftParen  TokenType = "LeftParen"
	RightParen TokenType = "RightParen"
	LeftBrace  TokenType = "LeftBrace"
	RightBrace TokenType = "RightBrace"
	Comma      TokenType = "Comma"
	Semicolon  TokenType = "Semicolon"
)

type Token struct {
	Type TokenType
	Val  string
}
