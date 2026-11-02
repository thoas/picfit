package failure

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
	httppkg "github.com/thoas/picfit/http"
)

type Handler func(*gin.Context) error

func Handle(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			if stderrors.Is(err, httppkg.ErrURLNotAllowed) {
				c.String(http.StatusBadRequest, "URL is not allowed")
				c.Abort()
				return
			}

			cerr := errors.Cause(err)

			if cerr == ErrFileNotExists || cerr == ErrKeyNotExists {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			if cerr == ErrFileNotModified {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}
			if cerr == ErrFileMaxDimensions {
				c.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}

			switch cerr.(type) {
			case binding.Errors:
				c.String(http.StatusBadRequest, cerr.Error())
				return
			}

			c.Error(err)
			panic(err)
		}
	}
}
