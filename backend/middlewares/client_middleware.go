package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var clientRoleValues = map[string]struct{}{
	"CLIENT":  {},
	"CLIENTE": {},
}

func ClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		role, ok := roleValue.(string)
		if !ok || strings.TrimSpace(role) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		if _, allowed := clientRoleValues[strings.ToUpper(strings.TrimSpace(role))]; !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "client access required"})
			return
		}

		c.Next()
	}
}
