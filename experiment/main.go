package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Node struct {
	Value string
	Left  *Node
	Right *Node
}

func NewNode(value string) *Node {
	return &Node{
		Value: value,
		Left:  nil,
		Right: nil,
	}
}

func PrintTree(n *Node, marker string, level int) {
	if n != nil {
		PrintTree(n.Right, "  "+marker, level+1)
		ident := strings.Repeat("  ", level)
		fmt.Printf("%s%s%s\n", ident, marker, n.Value)
		PrintTree(n.Left, "  "+marker, level+1)
	}
}

var idx = 0
var content = "45*3*(4+2)"
var contentLen = len(content)

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func terminal() *Node {
	res := ""
	for idx < contentLen && isNumeric(string(content[idx])) {
		res += string(content[idx]) //use string builder
		idx += 1
	}
	if res != "" {
		return NewNode(res)
	}
	if string(content[idx]) == "(" { //maybe we can use byte comparison here?
		idx += 1
		eval := eval()
		if idx < contentLen && string(content[idx]) == ")" {
			idx += 1
			return eval
		}
	}
	panic("failed to parse in terminal function")
}

func operation() *Node {
	res := terminal()

	for idx < contentLen && slices.Contains([]string{"*"}, string(content[idx])) {
		op := string(content[idx])
		idx += 1
		if op == "*" {
			n := NewNode("mul")
			n.Left = res
			n.Right = terminal()
			res = n
		} else {
			panic("failed in operation")
		}
	}
	return res
}

func eval() *Node {
	res := operation()
	for idx < contentLen && slices.Contains([]string{"+"}, string(content[idx])) {
		op := string(content[idx])
		idx += 1
		if op == "+" {
			n := NewNode("plus")
			n.Left = res
			n.Right = operation()
			res = n
		}
	}
	return res
}

func main() {
	fmt.Println("Math:   ", content)
	PrintTree(eval(), "--", 0)
}
