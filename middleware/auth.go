package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/picfit/util"
)

// Security wraps the request and confront sent parameters with secret key
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.FromContext(c)

		secretKey := cfg.SecretKey

		if secretKey != "" {
			params := url.Values{}

			for k, v := range util.SortMapString(c.MustGet("parameters").(map[string]string)) {
				params.Set(k, v)
			}

			if !signature.VerifySign(secretKey, params.Encode()) {
				c.String(http.StatusUnauthorized, "Invalid signature")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
