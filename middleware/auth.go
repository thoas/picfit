package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/signature"
)

// Security wraps the request and confront sent parameters with secret key
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.FromContext(c)

		secretKey := cfg.SecretKey

		if secretKey != "" {
			if !signature.VerifyParameters(secretKey, c.MustGet("parameters").(map[string]string)) {
				c.String(http.StatusUnauthorized, "Invalid signature")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
