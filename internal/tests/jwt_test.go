package tests

import (
	"os"
	"testing"
	"time"

	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateAndParseJWT(t *testing.T) {
	user := &models.User{
		ID:       3,
		Username: "admin2",
		Role:     "admin",
	}
	token, err := utils.GenerateJWT(user)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := utils.ParseJWT(token)
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	if claims["user_id"] != float64(user.ID) {
		t.Errorf("Expected user_id %v, got %v", user.ID, claims["user_id"])
	}
	if claims["username"] != user.Username {
		t.Errorf("Expected username %v, got %v", user.Username, claims["username"])
	}
	if claims["role"] != user.Role {
		t.Errorf("Expected role %v, got %v", user.Role, claims["role"])
	}
}

func TestParseJWTAllowExpired(t *testing.T) {
	user := &models.User{
		ID:       2,
		Username: "admin1",
		Role:     "admin",
	}
	// Create a token that expired 1 hour ago, so exp is in the past
	exp := time.Now().Add(-1 * time.Hour).Unix()
	claims := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("Failed to sign expired JWT: %v", err)
	}

	// fmt.Println("Generated expired token:", tokenString)
	// fmt.Printf("Claims used: %+v\n", claims)
	// fmt.Println("exp (as time):", time.Unix(exp, 0))
	// fmt.Println("now:", time.Now())

	// Should fail with ParseJWT and get err and will goes to the error block
	_, err = utils.ParseJWT(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, got none")
	} else {
		t.Logf("ParseJWT : %v", err)
	}

	// Should succeed with ParseJWTAllowExpired and err will be nil
	claimsMap, err := utils.ParseJWTAllowExpired(tokenString)
	if err != nil {
		t.Fatalf("ParseJWTAllowExpired failed: %v", err)
	}
	t.Logf("Claims returned by ParseJWTAllowExpired: %+v", claimsMap)
	if claimsMap["user_id"] != float64(user.ID) {
		t.Errorf("Expected user_id %v, got %v", user.ID, claimsMap["user_id"])
	}
}
