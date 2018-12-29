// Package main 提供一个CLI界面
package main

import (
	"fmt"
)

func test(a *[]string) {
	*a = []string{"asa", "aaa", "aa"}
}

func main() {
	a := []string{"1", "2", "3"}

	test(&a)
	fmt.Println(a)
}
