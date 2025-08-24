package main

import (
	"fmt"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d started\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d completed\n", id)
}

func main() {
	fmt.Println("Main started")

	// 启动协程
	go worker(1)
	go worker(2)

	fmt.Println("Main continues execution")
	time.Sleep(2 * time.Second) // 等待协程完成
	fmt.Println("Main finished")
}
