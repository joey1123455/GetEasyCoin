package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoveryWithFileLogger(t *testing.T) {
	// Setup
	logFilePath := "test_panic_log.txt"
	// defer os.Remove(logFilePath) // Clean up after test

	// Create a Gin router and add the middleware
	router := gin.New()
	router.Use(RecoveryWithFileLogger(logFilePath))

	// Add a route that will cause a panic
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// Simulate a request to the panic route
	req, _ := http.NewRequest("GET", "/panic", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the response status code
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", resp.Code)
	}

	// Check that the panic was logged to the file
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Error reading log file: %v", err)
	}

	if !strings.Contains(string(logContent), "Panic recovered: ") {
		t.Errorf("Expected log file to contain 'Panic recovered: ', got: %s", string(logContent))
	}
}
