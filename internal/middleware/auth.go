package middleware

import (
	"net/http"
	"strings"

	"lunch_menu/internal/database"
	"lunch_menu/internal/models"
	"lunch_menu/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// AuthMiddleware checks for a valid JWT token in the Authorization header.
func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Check if token is blacklisted
	if blacklisted, err := database.IsTokenBlacklisted(tokenString); err == nil && blacklisted {
		utils.Respond(c, http.StatusUnauthorized, "Unauthorized", nil, &models.ErrorResponse{
			Error:   "TOKEN_REVOKED",
			Message: "Token has been revoked",
		})
		c.Abort()
		return
	}

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		if err.Error() == "token is expired" || strings.Contains(err.Error(), "expired") {
			// Try to renew using refresh token from cookie
			refreshToken, cookieErr := c.Cookie("refresh_token")
			if cookieErr == nil && refreshToken != "" {
				newAccessToken, renewErr := utils.RenewAccessToken(refreshToken, tokenString)
				if renewErr == nil && newAccessToken != "" {
					// Optionally set new access token as cookie or header
					c.Header("X-New-Access-Token", newAccessToken)
					// Optionally parse new token and set claims
					newClaims, parseErr := utils.ParseJWT(newAccessToken)
					if parseErr == nil {
						c.Set("userClaims", newClaims)
						c.Next()
						return
					}
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
	// Access claims:
	// userID := claims["user_id"]
	// username := claims["username"]
	// role := claims["role"]
	// Store claims in context for use in handlers
	c.Set("userClaims", claims)
	// refreshToken, err := c.Cookie("refresh_token")
	// fmt.Printf("AuthMiddleware: refresh_token cookie=%q, err=%v\n", refreshToken, err)
	c.Next()
}

// AdminMiddleware checks if the user has admin role or will be used as role based authentication.
// we are not using it currently but can be used in future if needed.
func AdminMiddleware(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No user claims found"})
		return
	}
	userClaims, ok := claims.(map[string]interface{})
	if !ok || userClaims["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}
	c.Next()
}

// RateLimitMiddleware returns a Gin middleware that limits requests per IP.
// rateFormat example: "10-S" (10 requests per second), "100-M" (100 per minute)
func RateLimitMiddleware(rateFormat string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateFormat)
	if err != nil {
		panic("invalid rate format for rate limiter")
	}
	store := memory.NewStore()
	return ginlimiter.NewMiddleware(limiter.New(store, rate))
}
