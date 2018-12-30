// Package main 提供一个CLI界面
package main

import (
	"flag"
	"fmt"
	"prompt"
)

func main() {
	flag.Parse()
	flag.Usage()
	for _, v := range prompt.Strs {
		fmt.Println(*v)
	}

}
