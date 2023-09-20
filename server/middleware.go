package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func recoverMiddleware(c *gin.Context) {
	defer func() {
		// We must abort if there has been a recover, otherwise the
		// remaining handlers would be called. If everything went fine,
		// c.Abort is a no-op.
		c.Abort()
	}()

	defer func() {
		if rec := recover(); rec != nil {
			status := http.StatusInternalServerError
			http.Error(c.Writer, http.StatusText(status), status)
		}
	}()

	c.Next()
}
