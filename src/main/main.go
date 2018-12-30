// Package main 提供一个CLI界面
package main

import (
	"flag"
	"fmt"
	"prompt"
)

func main() {
	flag.Parse()
	if prompt.Help {
		flag.Usage()
	}

	fmt.Println("prompt.Help: ", prompt.Help)
	fmt.Println("prompt.IsMaster: ", prompt.IsMaster)
	fmt.Println("prompt.IsWorker: ", prompt.IsWorker)

	fmt.Println("prompt.MasterAddress: ", prompt.MasterAddress)
	fmt.Println("prompt.WorkerAddress: ", prompt.WorkerAddress)
	fmt.Println("prompt.Name: ", prompt.Name)
	fmt.Println("prompt.Filepath: ", prompt.Filepath)

}
