package main

import (
	"fmt"
	"os"
)

func main() {
	allArgs := os.Args
	if len(allArgs) < 2 {
		fmt.Printf("provide file to lex")
		return
	}
	fileContent, err := os.ReadFile(allArgs[1])
	if err != nil {
		panic(err)
	}
	rawData := string(fileContent)

	fmt.Println("Content: \n", rawData)

	lexRes := lex(rawData)
	// fmt.Println(lexRes)

	// AST logic

	ast := Ast{
		tokens:   lexRes,
		tokenLen: len(lexRes),
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
	fmt.Println("le start of tree")
	printTree(ast.root)
}
