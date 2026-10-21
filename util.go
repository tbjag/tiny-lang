package main

import "os"

func readFile(filename string) string {
	fileContent, _ := os.ReadFile(filename)
	return string(fileContent)
}