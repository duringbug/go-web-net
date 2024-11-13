package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/duringbug/go-web-net/pkg/logger" // 导入自定义的日志包
)

// Config 结构体表示配置文件中的配置
type Config struct {
	Server struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"server"`
}

// loadConfig 从指定的 JSON 文件加载配置，若加载失败则使用默认配置
func loadConfig(configPath string, log *logger.Logger) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		log.Warn("配置文件加载失败，使用默认配置")
		return &Config{
			Server: struct {
				Port int    `json:"port"`
				Host string `json:"host"`
			}{
				Port: 8080,
				Host: "localhost",
			},
		}, nil
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Error("解析配置文件失败: ", err)
		return nil, err
	}

	return &config, nil
}

// 处理接收到的 UDP 消息
func handleUDPMessage(conn *net.UDPConn, log *logger.Logger) {
	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Error("接收消息失败: ", err)
			continue
		}
		log.Info(fmt.Sprintf("接收到来自 %s 的消息: %s", addr, string(buffer[:n])))
	}
}

func main() {
	// 创建 logger 实例
	logFilePath := "log/cells.log"
	log, err := logger.NewLogger(logFilePath)
	if err != nil {
		log.Error("无法创建日志文件: ", err)
		return
	}
	defer log.Close()

	// 定义命令行参数
	configPath := flag.String("conf", "configs/config.json", "配置文件路径")
	flag.Parse()

	// 加载配置文件
	config, err := loadConfig(*configPath, log)
	if err != nil {
		log.Error("加载配置失败: ", err)
		return
	}

	// 输出加载的配置
	log.Info(fmt.Sprintf("服务器启动在 %s:%d", config.Server.Host, config.Server.Port))

	// 构建 UDP 服务器地址
	udpAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Info(fmt.Sprintf("启动 UDP 服务器，监听地址: %s", udpAddr))

	// 解析 UDP 地址
	address, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		log.Error("解析 UDP 地址失败: ", err)
		return
	}

	// 创建 UDP 连接
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Error("创建 UDP 连接失败: ", err)
		return
	}
	defer conn.Close()

	// 启动一个 goroutine 处理接收到的消息
	go handleUDPMessage(conn, log)

	// 启动 UDP 客户端发送消息（可以替换为其他逻辑）
	clientAddr := fmt.Sprintf("%s:%d", "localhost", config.Server.Port)
	remoteAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		log.Error("解析 UDP 客户端地址失败: ", err)
		return
	}

	// 发送测试消息
	message := []byte("Hello from UDP Client!")
	_, err = conn.WriteToUDP(message, remoteAddr)
	if err != nil {
		log.Error("发送消息失败: ", err)
		return
	}

	log.Info("消息已发送，等待接收...")
	// 阻塞以维持连接
	select {}
}
