// Package main 提供一个CLI界面
package main

import (
	"fmt"
)

type test struct {
	v1 string
	v2 string
}

func main() {
	var t *test
	t = &test{"asda", "aaa"}
	fmt.Println(t.v1)

}
