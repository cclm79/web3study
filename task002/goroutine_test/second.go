package goroutine_test

import (
	"fmt"
	"sync"
	"time"
)

// Task 表示一个可执行的任务
type Task func()

// TaskResult 存储任务执行结果
type TaskResult struct {
	Name     string
	Duration time.Duration
}

// TaskScheduler 并发任务调度器
type TaskScheduler struct {
	tasks map[string]Task
}

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks: make(map[string]Task),
	}
}

// AddTask 添加任务到调度器
func (ts *TaskScheduler) AddTask(name string, task Task) {
	ts.tasks[name] = task
}

// Run 并发执行所有任务并返回执行结果
func (ts *TaskScheduler) Run() map[string]time.Duration {
	var wg sync.WaitGroup
	results := make(chan TaskResult, len(ts.tasks))

	// 为每个任务启动协程
	for name, task := range ts.tasks {
		wg.Add(1)
		go func(name string, task Task) {

			start := time.Now()
			task() // 执行任务函数
			duration := time.Since(start)

			results <- TaskResult{Name: name, Duration: duration}

			defer wg.Done()

		}(name, task)
	}

	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	resultMap := make(map[string]time.Duration)
	for res := range results {
		resultMap[res.Name] = res.Duration
	}

	return resultMap
}

// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。
func TestSecond() {
	// 创建任务调度器
	scheduler := NewTaskScheduler()

	// 添加示例任务
	scheduler.AddTask("Task1", func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Task1 completed")
	})

	scheduler.AddTask("Task2", func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Task2 completed")
	})

	scheduler.AddTask("Task3", func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Task3 completed")
	})

	// 执行任务并获取结果
	results := scheduler.Run()

	// 打印执行时间统计
	fmt.Println("\nExecution times:")
	for name, duration := range results {
		fmt.Printf("%s: %v\n", name, duration)
	}
}
