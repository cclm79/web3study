package goroutine_test

import (
	"fmt"
	"time"
)

func printOdd() {
	for i := 1; i <= 10; {
		fmt.Println(i)
		i += 2
	}

}

func TestFirst() {
	fmt.Println("test start")
	go printOdd()

	go func() {
		for i := 0; i <= 10; {
			fmt.Println(i)
			i += 2
		}
	}()

	time.Sleep(10000 * time.Millisecond)
	fmt.Println("test end")
}
