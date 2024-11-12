package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

// 执行命令的函数
func runCommand(command string, args []string, wg *sync.WaitGroup) {
	// 结束时通知 WaitGroup
	defer wg.Done()

	// 准备命令
	cmd := exec.Command(command, args...)

	// 设置命令的输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		log.Printf("命令执行失败: %v", err)
	}
}

// main 函数，启动并行命令
func main() {
	var wg sync.WaitGroup

	// 需要执行的命令
	commands := []struct {
		command string
		args    []string
	}{
		{"./build/cell", []string{"-conf", "./configs/cells_config/cell_config01.json"}},
		{"./build/cell", []string{"-conf", "./configs/cells_config/cell_config02.json"}},
	}

	// 启动并行任务
	for _, cmd := range commands {
		wg.Add(1)                                 // 每个命令都需要通知 WaitGroup
		go runCommand(cmd.command, cmd.args, &wg) // 使用 goroutine 执行命令
	}

	// 等待所有命令执行完毕
	wg.Wait()

	fmt.Println("所有命令执行完毕")
}
