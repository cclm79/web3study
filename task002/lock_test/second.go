package lock_test

/**
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func TestSecond() {
	var counter int64 // 原子计数器
	var wg sync.WaitGroup

	// 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 每个协程递增计数器1000次
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait() // 等待所有协程完成
	fmt.Println("Counter value:", atomic.LoadInt64(&counter))
}
