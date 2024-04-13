package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// RecoveryWithFileLogger returns a Gin middleware function that recovers from panics and logs them to a file.
//
// logFilePath string: The file path where the panics will be logged.
// gin.HandlerFunc: The middleware function for Gin.
func RecoveryWithFileLogger(logFilePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic to a file
				writePanicToFile(logFilePath, err)

				// Respond with an internal server error
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

// writePanicToFile writes the panic message and stack trace to a specified log file.
//
// logFilePath is the path to the log file. err is the recovered panic message.
func writePanicToFile(logFilePath string, err interface{}) {
	// Open or create the log file in append mode
	file, fileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v\n", fileErr)
		return
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("Panic recovered: %s\n", time.Now().UTC().String()))
	file.Write(debug.Stack())
}
