package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"sync"

	"github.com/duringbug/go-web-net/pkg/logger" // 导入自定义的日志包
)

type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Config struct {
	Commands []Command `json:"commands"`
}

// loadConfig 从 JSON 文件加载配置
func loadConfig(filePath string) (*Config, error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// 创建 logger 实例
	logFilePath := "log/organsys.log"
	log, err := logger.NewLogger(logFilePath)
	if err != nil {
		log.Error("无法创建日志文件: ", err)
		return
	}
	defer log.Close()

	// 定义命令行参数 -conf，用于指定配置文件路径
	confPath := flag.String("conf", "configs/organsys_config/organ_config01.json", "配置文件的路径")
	flag.Parse() // 解析命令行参数

	// 加载配置文件
	config, err := loadConfig(*confPath)
	if err != nil {
		log.Error("加载配置失败: ", err)
		return
	}

	var wg sync.WaitGroup

	// 遍历并并行执行每个命令
	for _, cmd := range config.Commands {
		wg.Add(1) // 增加 WaitGroup 计数器
		go func(cmd Command) {
			defer wg.Done() // 命令执行完成后减小计数器
			log.Info("执行命令: ", cmd.Command, cmd.Args)

			// 使用 os/exec 包执行命令
			execCmd := exec.Command(cmd.Command, cmd.Args...)
			output, err := execCmd.CombinedOutput()
			if err != nil {
				log.Error("命令执行失败: ", err, " 命令: ", cmd.Command, cmd.Args)
			} else {
				log.Info("命令输出:\n", string(output))
			}
		}(cmd) // 这里传递 cmd 作为参数
	}

	// 等待所有 Goroutine 执行完成
	wg.Wait()
}
