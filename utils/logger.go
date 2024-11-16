package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(message string) {
	l.logger.SetPrefix("[INFO] ")
	l.logger.Println(message)
}

func (l *Logger) Error(message string) {
	l.logger.SetPrefix("[ERROR] ")
	l.logger.Println(message)
}

func (l *Logger) Warn(message string) {
	l.logger.SetPrefix("[WARN] ")
	l.logger.Println(message)
}

func GinRequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		log.Printf("[REQUEST] %s %s | Status: %d | Duration: %s | IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}
