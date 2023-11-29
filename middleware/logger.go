package middleware

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)


// Logger middleware function
func Logger() gin.HandlerFunc {
	logFilePath := "./logs/app.log" // Path to the log file

	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		if err := ensureLogDirExists(logFilePath); err != nil {
			log.Println("Error creating logs directory:", err)
			return
		}

		// Log request details to a file
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		log.Printf(
			"%s %s %s %v",
			c.Request.Method,
			c.Request.RequestURI,
			c.ClientIP(),
			time.Since(start),
		)
	}
}

func ensureLogDirExists(logFilePath string) error {
	dir := filepath.Dir(logFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
