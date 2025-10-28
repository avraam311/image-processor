package server

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/avraam311/image-processor/internal/api/handlers/images"
	"github.com/avraam311/image-processor/internal/middlewares"
)

func NewRouter(ginMode string, handlerIm *images.Handler) *ginext.Engine {
	e := ginext.New(ginMode)

	e.Use(middlewares.CORSMiddleware())
	e.Use(ginext.Logger())
	e.Use(ginext.Recovery())

	api := e.Group("/image-processor/api")
	{
		api.POST("/upload", handlerIm.UploadImage)
		api.GET("/image/:id", handlerIm.GetProcessedImage)
		api.DELETE("/image/:id", handlerIm.DeleteImage)
	}

	return e
}

func NewServer(addr string, router *ginext.Engine) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
