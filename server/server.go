package server

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thoas/picfit/config"
	"golang.org/x/net/context"
)

// Run loads a new server
func Run(ctx context.Context) {
	router := gin.Default()

	cfg := config.FromContext(ctx)

	// methods := map[string]views.View{
	// 	"redirect": views.RedirectView,
	// 	"display":  views.DisplayView,
	// 	"get":      views.GetView,
	// }
	//
	// for name, view := range methods {
	// 	router.GET(fmt.Sprintf("/%s", name), view)
	// 	router.GET(fmt.Sprintf("/%s/{sig}/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), view)
	// 	router.GET(fmt.Sprintf("/%s/{sig}/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", name), view)
	// 	router.GET(fmt.Sprintf("/%s/{sig}/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), view)
	// 	router.GET(fmt.Sprintf("/%s/{op}/x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), view)
	// 	router.GET(fmt.Sprintf("/%s/{op}/{w:[\\d]+}x/{path:[\\w\\-/.]+}", name), view)
	// 	router.GET(fmt.Sprintf("/%s/{op}/{w:[\\d]+}x{h:[\\d]+}/{path:[\\w\\-/.]+}", name), view)
	// }
	//
	// if cfg.Options.EnableUpload {
	// 	router.POST("/upload", views.UploadView)
	// }
	//
	// if cfg.Options.EnableDelete {
	// 	router.DELETE("/{path:[\\w\\-/.]+}", views.DeleteView)
	// }

	router.Run(fmt.Sprintf(":%s", strconv.Itoa(cfg.Port)))
}
