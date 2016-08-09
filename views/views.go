package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/payload"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/util"
)

const sigParamName = "sig"

// View represents a view interface
type View func(c *gin.Context)

// RequestKey injects an unique key from query parameters
func RequestKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryString := make(map[string]string)

		for _, param := range c.Params {
			queryString[param.Key] = param.Value
		}

		for k, v := range c.Request.URL.Query() {
			queryString[k] = v[0]
		}

		sorted := util.SortMapString(queryString)

		delete(sorted, sigParamName)

		serialized := hash.Serialize(sorted)

		key := hash.Tokey(serialized)

		c.Set("key", key)

		c.Next()
	}
}

// DisplayView displays and image using resizing parameters
func DisplayView(c *gin.Context) {
	file, err := application.ImageFileFromRequest(c, true, true)

	if err != nil {
		panic(err)
	}

	for k, v := range file.Headers {
		c.Header(k, v)
	}

	c.Data(http.StatusOK, "", file.Content())
}

// UploadView uploads an image to the destination storage
func UploadView(c *gin.Context) {
	multipartPayload := new(payload.MultipartPayload)
	errs := binding.Bind(c.Request, multipartPayload)
	if errs.Handle(c.Writer) {
		return
	}

	file, err := multipartPayload.Upload(storage.SourceFromContext(c))

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})
}

// DeleteView deletes a file from storages
func DeleteView(c *gin.Context) {
	err := application.ImageCleanup(c, c.Param("path"))

	if err != nil {
		panic(err)
	}

	c.String(http.StatusOK, "Ok")
}

// GetView generates an image synchronously and return its information from storages
func GetView(c *gin.Context) {
	file, err := application.ImageFileFromRequest(c, false, false)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})
}

// RedirectView redirects to the image using base url from storage
func RedirectView(c *gin.Context) {
	file, err := application.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusMovedPermanently, file.URL())
}
