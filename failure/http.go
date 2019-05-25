package failure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/pkg/errors"
)

type Handler func(*gin.Context) error

func Handle(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)
		if err != nil {
			cerr := errors.Cause(err)

			if cerr == ErrFileNotExists || cerr == ErrKeyNotExists {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			if cerr == ErrFileNotModified {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}

			switch cerr.(type) {
			case binding.Errors:
				c.String(http.StatusBadRequest, cerr.Error())
			}

			panic(err)
		}
	}
}
