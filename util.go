package main

import "os"

func readFile(filename string) string {
	fileContent, _ := os.ReadFile(filename)
	return string(fileContent)
}

func charCode(lit string) int {
	body := lit[1 : len(lit)-1]
	switch body {
	case `\n`:
		return '\n'
	case `\\`:
		return '\\'
	}
	return int(body[0])
}
