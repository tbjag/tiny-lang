package main

import (
	"testing"
)

func TestAst(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"tests/ast_else.txt", `Sequence(nil, If(Less(p, 3), If(Sequence(nil, Assign(p, Add(p, 1))), Sequence(nil, Assign(x, Add(x, 3))))))`},
		{"tests/ast_mixed_precedence.txt", `Sequence(nil, Assign(x, Subtract(Add(1, Multiply(2, 3)), Mod(Divide(4, 2), 3))))`},
		{"tests/prime_number_generation.txt", `Sequence(Sequence(Sequence(Sequence(Sequence(nil, Assign(count, 1)), Assign(n, 1)), Assign(limit, 100)), While(Less(n, limit), Sequence(Sequence(Sequence(Sequence(Sequence(nil, Assign(k, 3)), Assign(p, 1)), Assign(n, Add(n, 2))), While(And(LessEqual(Multiply(k, k), n), p), Sequence(Sequence(nil, Assign(p, NotEqual(Multiply(Divide(n, k), k), n))), Assign(k, Add(k, 2))))), If(p, If(Sequence(Sequence(nil, Sequence(Sequence(nil, Prti(n, nil)), Prts(" is prime\n", nil))), Assign(count, Add(count, 1))), nil))))), Sequence(Sequence(Sequence(nil, Prts("Total primes found: ", nil)), Prti(count, nil)), Prts("\n", nil)))`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := readFile(tt.name)
			res := printTreeNested(createAst(lex(content)))
			if res != string(tt.want) {
				t.Errorf("Parsed AST does not match result:\n   %s\n-- %s\nin file %s", res, tt.want, tt.name)
			}
		})
	}
}
