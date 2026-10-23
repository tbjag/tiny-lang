package main

import (
	"strings"
	"testing"
)

func TestEmit(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"tests/count.txt", "validation/count.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emit := gen(createAst(lex(readFile(tt.name))))
			want := readFile(tt.want)
			if strings.TrimSpace(emit) != strings.TrimSpace(want) {
				t.Errorf("Emitted AST does not match result:\n for file %s", tt.want)
			}
		})
	}
}
