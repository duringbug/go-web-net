package logger

import (
	"io"
	"log"
	"os"
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
	l.infoLogger.Println(v...)
}

// Warn 打印警告日志
func (l *Logger) Warn(v ...interface{}) {
	l.warnLogger.Println(v...)
}

// Error 打印错误日志
func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}
