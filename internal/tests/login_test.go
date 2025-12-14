package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"lunch_menu/internal/config"
	"lunch_menu/internal/database"
	"lunch_menu/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TestMain sets up the config and database connection for all tests in this package.
func TestMain(m *testing.M) {

	_ = godotenv.Load("../../.env.test") // 1. Load env vars from file
	config.MustLoadConfig()              // 2. Populate config.AppConfig

	// Initialize database connection
	if err := database.InitDatabase(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	// Run tests
	code := m.Run()
	// Cleanup
	_ = database.CloseDatabase()
	os.Exit(code)
}

func TestAdminLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/user/login", handlers.UserLogin)

	// Prepare request body
	body := map[string]string{
		"username":      "admin1",
		"password_hash": "admin1", // Use correct password for your test DB
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent") // <-- Add this line
	req.RemoteAddr = "127.0.0.1:12345"         // <-- Add this line

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	// Optionally, check for Set-Cookie headers, response body, etc.
}
