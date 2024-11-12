package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// Config 结构体表示配置文件中的配置
type Config struct {
	Server struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"server"`
}

// loadConfig 从指定的 JSON 文件加载配置，若加载失败则使用默认配置
func loadConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		log.Println("配置文件加载失败，使用默认配置")
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
		return nil, err
	}

	return &config, nil
}

// 处理接收到的 UDP 消息
func handleUDPMessage(conn *net.UDPConn) {
	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("接收消息失败: %v", err)
			continue
		}
		log.Printf("接收到来自 %s 的消息: %s", addr, string(buffer[:n]))
	}
}

func main() {
	// 定义命令行参数
	configPath := flag.String("conf", "configs/config.json", "配置文件路径")
	flag.Parse()

	// 加载配置文件
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 输出加载的配置
	log.Printf("服务器启动在 %s:%d", config.Server.Host, config.Server.Port)

	// 构建 UDP 服务器地址
	udpAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("启动 UDP 服务器，监听地址: %s", udpAddr)

	// 解析 UDP 地址
	address, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		log.Fatalf("解析 UDP 地址失败: %v", err)
	}

	// 创建 UDP 连接
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalf("创建 UDP 连接失败: %v", err)
	}
	defer conn.Close()

	// 启动一个 goroutine 处理接收到的消息
	go handleUDPMessage(conn)

	// 启动 UDP 客户端发送消息（可以替换为其他逻辑）
	clientAddr := fmt.Sprintf("%s:%d", "localhost", config.Server.Port)
	remoteAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		log.Fatalf("解析 UDP 客户端地址失败: %v", err)
	}

	// 发送测试消息
	message := []byte("Hello from UDP Client!")
	_, err = conn.WriteToUDP(message, remoteAddr)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	log.Println("消息已发送，等待接收...")
	// 阻塞以维持连接
	select {}
}
