// Package main 提供一个CLI界面
package main

import "fmt"

func test(a *[]string) {
	*a = []string{"asa", "aaa", "aa"}
}

func main() {
	a := "1:345"
	fmt.Println(a[:-1], a[1+1:])
}
