package main

import (
	"os"
	"testing"
)

func readFile(filename string) string {
	fileContent, _ := os.ReadFile(filename)
	return string(fileContent)
}

func TestBasicFiles(t *testing.T) {
	tests := []struct {
		name string
		want []TokenType
	}{
		{"tests/hello_world.txt", []TokenType{
			Print, LeftParen, String,
			RightParen, Semicolon, EOF}},
		{"tests/variables.txt", []TokenType{
			Identifier, Assign, Integer,
			Semicolon, Print, LeftParen,
			Identifier, Comma, String,
			RightParen, Semicolon, EOF}},
		{"tests/all_lex.txt", []TokenType{
			Print, Subtract, Put,
			Less, If, Greater,
			Else, LessEqual, While,
			GreaterEqual, LeftBrace, Equal,
			RightBrace, NotEqual, LeftParen,
			And, RightParen, Or,
			Subtract, Semicolon, Not,
			Comma, Multiply, Assign,
			Divide, Integer, Mod,
			String, Add, Identifier,
			Character, Character, Character,
			EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := readFile(tt.name)
			res := lex(content)
			if len(res) != len(tt.want) {
				t.Errorf("Parsed token lengths to not match, got %d want %d", len(res), len(tt.want))
			}

			for i := range len(res) {
				if res[i].Type != tt.want[i] {
					t.Errorf("Parsed tokens to not match, got %s want %s", res[i].Type, tt.want[i])
				}
			}
		})
	}
}
