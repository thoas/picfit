package server

import "github.com/gin-gonic/gin"

type handlerMethod func(string, ...gin.HandlerFunc) gin.IRoutes
type endpoint struct {
	handler gin.HandlerFunc
	method  handlerMethod
	pattern string
	route   string
}
