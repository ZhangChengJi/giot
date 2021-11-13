package handler

import (
	"github.com/gin-gonic/gin"
)

type RouteRegister interface {
	ApplyRoute(r *gin.Engine)
}
