package main

import (
	"fmt"
	"io"
	"runtime"
)

// func main() {
// 	var c1, c2, c3 chan int
// 	var i1, i2 int
// 	select {
// 	case i1 = <-c1:
// 		fmt.Printf("received ", i1, " from c1\n")
// 	case c2 <- i2:
// 		fmt.Printf("sent ", i2, " to c2\n")
// 	case i3, ok := (<-c3): // same as: i3, ok := <-c3
// 		if ok {
// 			fmt.Printf("received ", i3, " from c3\n")
// 		} else {
// 			fmt.Printf("c3 is closed\n")
// 		}
// 	default:
// 		fmt.Printf("no communication\n")
// 	}
// }
func createCounter(start int) chan int {
	next := make(chan int)
	go func(i int) {
		for {
			next <- i
			i++
		}
	}(start)
	return next
}

func main() {
	counterA := createCounter(2)
	counterB := createCounter(102)
	for i := 0; i < 5; i++ {
		a := <-counterA
		fmt.Printf("(A->%d, B->%d)\n", a, <-counterB)
	}
	var w io.Writer
	runtime.Gosched()
}
