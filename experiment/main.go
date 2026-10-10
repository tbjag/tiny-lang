package main

import (
	"fmt"
	"slices"
	"strconv"
)


var idx = 0
var content = "45*5+(333/3)+10"
var contentLen =  len(content)

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func terminal() int64 { // can we convert to reg. int?
	res := ""
	for idx < contentLen && isNumeric(string(content[idx])) {
		res += string(content[idx]) //use string builder
		idx += 1
	}
	if res != "" {
		num, err := strconv.ParseInt(res, 10, 64)
		if err != nil {
			panic("failed to parse int")
		}
		return num
	}
	if string(content[idx]) == "(" { //maybe we can use byte comparison here?
		idx += 1
		res2 := eval()
		if idx < contentLen && string(content[idx]) == ")" {
			idx += 1
			return int64(res2)
		}
	}
	panic("failed to parse in terminal function")
}

func operation() int64 {
	res := terminal()

	for idx < contentLen && slices.Contains([]string{"*", "/"}, string(content[idx])) {
		if string(content[idx]) == "*" {
			idx += 1
			res *= terminal()
		} else if string(content[idx]) == "/" {
			idx += 1
			res /= terminal()
		} else {
			panic("failed to parse in operation")
		}
	}

	return res
}

func eval() int64 {
	res := operation()

	for idx < contentLen && slices.Contains([]string{"+"}, string(content[idx])) {
		idx += 1
		res += operation()
	}
	
	return res
}

func main() {
	fmt.Println("Math:   ", content)
	fmt.Println("Result: ", eval())
}
