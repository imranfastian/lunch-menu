package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"lunch_menu/internal/database"
	"lunch_menu/internal/models"

	"github.com/golang-jwt/jwt/v4"
)

// ParseJWT parses and validates a JWT token string and returns the claims.
func ParseJWT(tokenString string) (map[string]interface{}, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	// Handle specific JWT errors
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			}
			if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.New("invalid signature")
			}
			return nil, fmt.Errorf("token validation error: %w", err)
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Explicitly check expiration (for tokens without standard validation)
	if exp, ok := claims["exp"].(float64); ok {
		now := time.Now().Unix()
		if int64(exp) < now {
			return nil, errors.New("token is expired")
		}
	}

	// Convert jwt.MapClaims to map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}
	return result, nil
}

// GenerateJWT generates a JWT token for the given user.
func GenerateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret" // For development only; use env in production!
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// RenewAccessToken tries to renew an access token using a valid refresh token.
// Returns new access token string, or error if not possible.
func RenewAccessToken(refreshToken string, expiredToken string) (string, error) {
	if refreshToken == "" || expiredToken == "" {
		return "", errors.New("missing refresh or access token")
	}

	claims, err := ParseJWTAllowExpired(expiredToken)
	if err != nil {
		return "", errors.New("could not parse expired access token")
	}

	userIDRaw, ok := claims["user_id"]
	if !ok {
		return "", errors.New("user_id missing in token")
	}
	var userID int
	switch v := userIDRaw.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		return "", errors.New("invalid user_id type in token")
	}

	rt, err := database.GetRefreshToken(userID, refreshToken)
	if err != nil || rt == nil {
		return "", errors.New("invalid or expired refresh token")
	}
	user, err := database.GetUserByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}
	return GenerateJWT(user)
}

// ParseJWTAllowExpired parses a JWT and returns claims even if the token is expired.
func ParseJWTAllowExpired(tokenString string) (map[string]interface{}, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	parser := jwt.Parser{
		SkipClaimsValidation: true, // disables exp/nbf checks
	}

	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Convert jwt.MapClaims to map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}
	return result, nil
}
