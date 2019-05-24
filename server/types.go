package server

import "github.com/gin-gonic/gin"

type handlerMethod func(string, ...gin.HandlerFunc) gin.IRoutes
type endpoint struct {
	pattern string
	handler gin.HandlerFunc
	method  handlerMethod
}
