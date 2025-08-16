package channel_test

import (
	"fmt"
	"sync"
)

func TestSecond() {

	// 创建缓冲通道，容量为10
	ch := make(chan int, 10)
	var wg sync.WaitGroup

	// 生产者协程：发送100个整数
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch) // 发送完成后关闭通道
		for i := 0; i < 100; i++ {
			fmt.Println("send:", i)
			ch <- i
		}
	}()

	// 消费者协程：接收并打印整数
	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range ch { // 自动检测通道关闭
			fmt.Println("Received:", num)
		}
	}()

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("All done!")

}
