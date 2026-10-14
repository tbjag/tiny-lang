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
	rootNode := Token{
		Type: EOF,
		Val: "",
	}
	ast := Ast {
		tokens: lexRes,
		tokenLen: len(lexRes),
		root: newNode(rootNode),
		idx: 0,
	}
	test := ast.expression()
	printTree(test, "---", 0)
}
