package main

import (
	"regexp"
)

type TokenMatch struct {
	TokenType TokenType
	Regex     *regexp.Regexp
	Size      int
}

func (t TokenMatch) ParseKeyword(data string, idx int) (int, Token, bool) {
	loc := t.Regex.FindStringIndex(data[idx:])
	token := Token{
		Type: t.TokenType,
		Val:  "",
	}
	if len(loc) != 2 || loc[0] != 0 {
		return idx, token, false
	}
	return idx + loc[1], token, true
}

func (t TokenMatch) ParseVariable(data string, idx int) (int, Token, bool) {
	loc := t.Regex.FindStringIndex(data[idx:])
	token := Token{
		Type: t.TokenType,
		Val:  "",
	}
	if len(loc) != 2 || loc[0] != 0 {
		return idx, token, false
	}
	token.Val = data[idx : idx+loc[1]]
	return idx + loc[1], token, true
}

var JunkMatches = []*regexp.Regexp{
	regexp.MustCompile(`\s+`),
	regexp.MustCompile(`\/\*[\s\S]*?\*\/`),
}

var KeywordTokenMatches = []TokenMatch{
	{Print, regexp.MustCompile(`\bprint\b`), 5},
	{Put, regexp.MustCompile(`\bputc\b`), 3},
	{While, regexp.MustCompile(`\bwhile\b`), 5},
	{If, regexp.MustCompile(`if`), 2},
	{Else, regexp.MustCompile(`else`), 4},
	{Equal, regexp.MustCompile(`==`), 2},
	{NotEqual, regexp.MustCompile(`!=`), 2},
	{LessEqual, regexp.MustCompile(`<=`), 2},
	{GreaterEqual, regexp.MustCompile(`>=`), 2},
	{And, regexp.MustCompile(`&&`), 2},
	{Or, regexp.MustCompile(`\|\|`), 2},
	{Greater, regexp.MustCompile(`>`), 1},
	{Less, regexp.MustCompile(`<`), 1},
	{Not, regexp.MustCompile(`!`), 2},
	{Multiply, regexp.MustCompile(`\*`), 1},
	{Multiply, regexp.MustCompile(`\*`), 1},
	{Divide, regexp.MustCompile(`/`), 1},
	{Mod, regexp.MustCompile(`%`), 1},
	{Add, regexp.MustCompile(`\+`), 1},
	{Subtract, regexp.MustCompile(`-`), 1},
	{Semicolon, regexp.MustCompile(`;`), 1},
	{Assign, regexp.MustCompile(`=`), 1},
	{Comma, regexp.MustCompile(`,`), 1},
	{LeftParen, regexp.MustCompile(`\(`), 1},
	{RightParen, regexp.MustCompile(`\)`), 1},
	{LeftBrace, regexp.MustCompile(`\{`), 1},
	{RightBrace, regexp.MustCompile(`\}`), 1},
}

var VariableTokenMatches = []TokenMatch{
	{String, regexp.MustCompile(`"[^"\\]*(\\.[^"\\]*)*"`), -1},
	{Integer, regexp.MustCompile(`[0-9]+`), -1},
	{Character, regexp.MustCompile(`'([^'\n]|\\n|\\\\)'`), -1},
	{Identifier, regexp.MustCompile(`[_a-zA-Z][_a-zA-Z0-9]*`), -1},
}

func checkForKeyword(data string, idx int) (int, Token, bool) {
	match := false
	token := Token{}
	for _, tm := range KeywordTokenMatches {
		if _idx, _token, ok := tm.ParseKeyword(data, idx); ok {
			match = true
			idx = _idx
			token = _token
			break
		}
	}
	return idx, token, match
}

func checkForVariable(data string, idx int) (int, Token, bool) {
	match := false
	token := Token{}
	for _, tm := range VariableTokenMatches {
		if _idx, _token, ok := tm.ParseVariable(data, idx); ok {
			match = true
			idx = _idx
			token = _token
			break
		}
	}
	return idx, token, match
}

func eatJunk(data string, idx int) (int, bool) {
	for _, jm := range JunkMatches {
		loc := jm.FindStringIndex(data[idx:])
		if len(loc) != 2 || loc[0] != 0 {
			continue
		}
		return idx + loc[1], true
	}
	return idx, false
}

func lex(content string) []Token {
	var idx int
	var res []Token
	contentLen := len(content)

	for {
		match := false

		if _idx, ok := eatJunk(content, idx); ok {
			idx = _idx
			continue
		}

		if _idx, token, ok := checkForKeyword(content, idx); ok {
			match = true
			idx = _idx
			res = append(res, token)
			continue
		}

		if _idx, token, ok := checkForVariable(content, idx); ok {
			match = true
			idx = _idx
			res = append(res, token)
			continue
		}

		if idx >= contentLen {
			res = append(res, Token{EOF, ""})
			break
		}

		if !match {
			panic("Could not lex: " + content[idx:])
		}

	}

	return res
}
