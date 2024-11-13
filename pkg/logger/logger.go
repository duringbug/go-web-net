package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

// NewLogger 创建一个新的 Logger 实例
func NewLogger(logFilePath string) (*Logger, error) {
	// 打开日志文件，支持追加写入
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// 创建带有多重输出的日志（控制台 + 文件）
	multiWriter := io.MultiWriter(file, os.Stdout)

	return &Logger{
		infoLogger:  log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		file:        file,
	}, nil
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// Info 打印信息日志
func (l *Logger) Info(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)                        // 获取文件名和行号
	logMessage := append(v, fmt.Sprintf(" [%s:%d]", file, line)) // 格式化为字符串
	l.infoLogger.Println(logMessage...)
}

// Warn 打印警告日志
func (l *Logger) Warn(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logMessage := append(v, fmt.Sprintf(" [%s:%d]", file, line))
	l.warnLogger.Println(logMessage...)
}

// Error 打印错误日志
func (l *Logger) Error(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logMessage := append(v, fmt.Sprintf(" [%s:%d]", file, line))
	l.errorLogger.Println(logMessage...)
}
