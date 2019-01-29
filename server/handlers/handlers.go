package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mholt/binding"

	api "gopkg.in/fukata/golang-stats-api-handler.v1"

	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/constants"
	"github.com/thoas/picfit/errs"
	"github.com/thoas/picfit/payload"
	"github.com/thoas/picfit/storage"
)

func StatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, api.GetStats())
}

// Healthcheck displays an ok response for healthcheck
func Healthcheck(uptime time.Time) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"uptime":     uptime,
			"status":     "Ok",
			"version":    constants.Version,
			"revision":   constants.Revision,
			"build_time": constants.BuildTime,
			"compiler":   constants.Compiler,
		})
	}
}

// Display displays and image using resizing parameters
func Display(c *gin.Context) {
	file, err := application.ImageFileFromContext(c, true, true)

	if err != nil {
		errs.Handle(err, c.Writer)

		return
	}

	for k, v := range file.Headers {
		c.Header(k, v)
	}

	c.Data(http.StatusOK, file.ContentType(), file.Content())
}

// Upload uploads an image to the destination storage
func Upload(c *gin.Context) {
	multipartPayload := new(payload.MultipartPayload)
	errs := binding.Bind(c.Request, multipartPayload)
	if errs != nil {
		c.String(http.StatusBadRequest, errs.Error())
		return
	}

	file, err := multipartPayload.Upload(storage.DestinationFromContext(c))

	if err != nil {
		c.String(http.StatusBadRequest, errs.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})
}

// Delete deletes a file from storages
func Delete(c *gin.Context) {
	err := application.Delete(c, c.Param("path")[1:])

	if err != nil {
		errs.Handle(err, c.Writer)

		return
	}

	c.String(http.StatusOK, "Ok")
}

// Get generates an image synchronously and return its information from storages
func Get(c *gin.Context) {
	file, err := application.ImageFileFromContext(c, false, false)

	if err != nil {
		errs.Handle(err, c.Writer)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})
}

// Redirect redirects to the image using base url from storage
func Redirect(c *gin.Context) {
	file, err := application.ImageFileFromContext(c, false, false)

	if err != nil {
		errs.Handle(err, c.Writer)

		return
	}

	c.Redirect(http.StatusMovedPermanently, file.URL())
}
