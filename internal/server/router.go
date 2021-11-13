package server

import (
	"giot/internal/handler"
	"giot/internal/handler/device"
	"giot/internal/log"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func routers() *gin.Engine {
	var Router = gin.New()
	// 跨域
	//Router.Use(middleware.Cors()) // 如需跨域可以打开
	log.Info("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Info("register swagger handler")
	//	b := &PublicGroup{Router.Group("")}
	//b.Register()
	p := &PrivateGroup{Router}
	p.Register()
	log.Info("router register success")
	return Router

}

type PublicGroup struct {
	*gin.Engine
}
type PrivateGroup struct {
	*gin.Engine
}

func (r *PrivateGroup) Register() {
	//r.Use(plugin.JwtAuth()) //认证
	list := []handler.RouteRegister{
		device.NewHandler(),
	}
	for _, register := range list {
		register.ApplyRoute(r.Engine)
	}
}

//func (r *PublicGroup) Register() {
//	list := []core.RouteRegister{
//		system.NewHandler(),
//	}
//	for _, register := range list {
//		register.Router(r.RouterGroup)
//	}
//}
