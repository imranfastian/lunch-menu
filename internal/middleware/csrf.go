package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFProtection checks for a valid CSRF token in the request header.
// For APIs used only by non-browser clients (mobile apps, Postman, etc.), CSRF is not typically needed.
// For browser-based apps, always use CSRF protection on state-changing endpoints.
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only check CSRF for state-changing methods
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch ||
			c.Request.Method == http.MethodDelete {

			csrfToken := c.GetHeader("X-CSRF-Token")
			sessionToken, err := c.Cookie("csrf_token")
			if err != nil || csrfToken == "" || csrfToken != sessionToken {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid or missing CSRF token"})
				return
			}
		}
		c.Next()
	}
}
