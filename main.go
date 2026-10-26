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

	fmt.Println("Content:\n", rawData)

	lexRes := lex(rawData)
	// fmt.Println(lexRes)
	ast := createAst(lexRes)
	// printTreePreOrder(ast)
	// fmt.Println(printTreeNested(ast))
	// fmt.Println("-----------------")
	g, instructions := gen(ast)
	fmt.Println(instructions)
	vm := Vm{instructions: g.instructions}
	vm.parseInstructions()
}
