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
	fmt.Println(lexRes)

	// AST logic
	rootNode := Token{
		Type: Semicolon,
		Val:  "root",
	}
	ast := Ast{
		tokens:   lexRes,
		tokenLen: len(lexRes),
		root:     newNode(rootNode),
		idx:      0,
	}
	for ast.tokens[ast.idx].Type != EOF {
		t := Token{
			Type: Semicolon,
			Val: "",
		}
		n := newNode(t)
		n.Left = ast.root
		n.Right = ast.stmt()
		ast.root = n
	}
	printTree(ast.root, "---", 0)
}
