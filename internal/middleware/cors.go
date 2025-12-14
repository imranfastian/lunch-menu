package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware sets CORS headers for API security and cross-origin support
// it is only browser issue, such issues never happen in postman or curl as these are server to server requests
// if credentials are involved, do not use "*" in Access-Control-Allow-Origin, instead specify the exact domain
// for production, replace "*" with your frontend domain
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Change "*" to your frontend domain in production
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // preflighted requests

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
