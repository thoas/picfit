package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thoas/go-funk"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/signature"
)

// Security wraps the request and confront sent parameters with secret key
func Security(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secretKey != "" {
			if !signature.VerifyParameters(secretKey, c.MustGet("parameters").(map[string]interface{})) {
				c.String(http.StatusUnauthorized, "Invalid signature")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func RestrictIPAddresses(ipAddresses []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(ipAddresses) > 0 {
			if !funk.InStrings(ipAddresses, c.ClientIP()) {
				c.String(http.StatusUnauthorized, "Endpoint restricted")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func RestrictSizes(sizes []config.AllowedSize) gin.HandlerFunc {
	handler := func(c *gin.Context, sizes []config.AllowedSize) {
		params := c.MustGet("parameters").(map[string]interface{})

		var w int
		var h int
		var err error

		if w, err = strconv.Atoi(params["w"].(string)); err != nil {
			return
		}

		if h, err = strconv.Atoi(params["h"].(string)); err != nil {
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
		if len(sizes) > 0 {
			handler(c, sizes)
		}

		c.Next()
	}
}
