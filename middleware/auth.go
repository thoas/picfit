package middleware

import (
	"net/http"
	"strconv"

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

func RestrictSizes() gin.HandlerFunc {
	handler := func(c *gin.Context, sizes []config.AllowedSize) {
		params := c.MustGet("parameters").(map[string]string)

		var w int
		var h int
		var err error

		if w, err = strconv.Atoi(params["w"]); err != nil {
			return
		}

		if h, err = strconv.Atoi(params["h"]); err != nil {
			return
		}

		ok := false
		for _, size := range sizes {
			if size.Height == h && size.Width == w {
				ok = true
				break
			}
		}

		if !ok {
			c.String(http.StatusForbidden, "Requested size not allowed")
			c.Abort()
		}
	}

	return func(c *gin.Context) {
		sizes := config.FromContext(c).Options.AllowedSizes

		if len(sizes) > 0 {
			handler(c, sizes)
		}

		c.Next()
	}
}
