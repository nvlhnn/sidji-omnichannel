package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sidji-omnichannel/internal/config"
)

// CORS middleware for handling cross-origin requests
func CORS(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Allow the configured frontend URL, development mode, or localhost
		allowedOrigin := cfg.FrontendURL
		isAllowed := origin == allowedOrigin || cfg.Env == "development" || strings.HasPrefix(origin, "http://localhost")
		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
