package server

import (
	"giot/internal/conf"
	"giot/internal/filter"
	"giot/internal/handler"
	"giot/internal/handler/device"
	"giot/internal/log"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func setupRouter() *gin.Engine {
	if conf.ENV == conf.EnvLOCAL || conf.ENV == conf.EnvDEV {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	logger := log.GetLogger(log.AccessLog)
	r.Use(filter.CORS(), filter.RequestId(), filter.IPFilter(), filter.RequestLogHandler(logger)) //filter.SchemaCheck(), filter.RecoverHandler())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	list := []handler.RouteRegister{
		device.NewHandler(),
	}
	for _, register := range list {
		register.ApplyRoute(r)
	}
	return r
}
