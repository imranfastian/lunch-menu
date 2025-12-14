package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"lunch_menu/internal/database"
	"lunch_menu/internal/middleware"
	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_RenewAccessToken(t *testing.T) {
	// Setup test user
	user := &models.User{
		ID:       2,
		Username: "admin1",
		Role:     "admin",
	}
	t.Logf("Test user: %+v", user)

	// Create an expired access token
	expiredExp := time.Now().Add(-1 * time.Hour).Unix()
	claims := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiredExp,
	}
	token := GenerateTestJWT(claims, t)
	// t.Logf("Generated expired access token: %s", token)
	// t.Logf("Claims: %+v", claims)
	// t.Logf("exp (as time): %v, now: %v", time.Unix(expiredExp, 0), time.Now())

	// Create and store a valid refresh token
	refreshToken := utils.GenerateRandomToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	t.Logf("Generated refresh token: %s (expires at %v)", refreshToken, expiresAt)
	err := database.SaveRefreshToken(user.ID, refreshToken, "test-agent", "127.0.0.1", expiresAt)
	if err != nil {
		t.Fatalf("Failed to save refresh token: %v", err)
	}

	// Setup Gin router with AuthMiddleware and a test handler
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with expired access token and refresh token cookie
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:12345" // Gin uses this for c.ClientIP()
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// t.Logf("Response code: %d", w.Code)
	// t.Logf("Response body: %s", w.Body.String())
	// t.Logf("X-New-Access-Token header: %s", w.Header().Get("X-New-Access-Token"))

	// Check response
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	if w.Header().Get("X-New-Access-Token") == "" {
		t.Error("Expected X-New-Access-Token header to be set")
	}
}

// Helper to generate a JWT with custom claims for testing
// (You can put this in internal/utils/jwt.go or in your test file)
func GenerateTestJWT(claims map[string]interface{}, t *testing.T) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign test JWT: %v", err)
	}
	return tokenString
}
