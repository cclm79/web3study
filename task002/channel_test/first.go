package channel_test

/**
题目 ：编写一个程序，使用通道实现两个协程之间的通信。
一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
*/

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// 只接收channel的函数
func receiveOnly(ch <-chan int) {
	defer wg.Done()
	for v := range ch {

		fmt.Printf("接收到: %d\n", v)
		time.Sleep(100 * time.Millisecond)
	}

}

// 只发送channel的函数
func sendOnly(ch chan<- int) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}
	close(ch)

}

func TestFirst() {
	// 创建一个带缓冲的channel
	ch := make(chan int, 10)

	wg.Add(2)

	// 启动发送goroutine
	go sendOnly(ch)

	// 启动接收goroutine
	go receiveOnly(ch)

	wg.Wait()

}
