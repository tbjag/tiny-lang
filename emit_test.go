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
		{"tests/hello_world.txt", "validation/hello_world.txt"},
		{"tests/variables.txt", "validation/phoenix_number.txt"},
		{"tests/100_doors.txt", "validation/100_doors.txt"},
		{"tests/deep.txt", "validation/deep.txt"},
		{"tests/factorial.txt", "validation/factorial.txt"},
		{"tests/fibonacci.txt", "validation/fibonacci.txt"},
		{"tests/fizz_buzz.txt", "validation/fizz_buzz.txt"},
		{"tests/greatest_divisor.txt", "validation/greatest_divisor.txt"},
		{"tests/mandelbrot.txt", "validation/mandelbrot.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emit := gen(createAst(lex(readFile(tt.name))))
			want := readFile(tt.want)
			if strings.TrimSpace(emit) != strings.TrimSpace(want) {
				t.Errorf("Emitted code does not match result:\n for file %s", tt.want)
			}
		})
	}
}
