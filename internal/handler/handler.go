package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"
	"net/http"
)

type RouteRegister interface {
	ApplyRoute(r *gin.Engine)
}

type ErrorTransformMiddleware struct {
	middleware.BaseMiddleware
}

func (mw *ErrorTransformMiddleware) Handle(ctx droplet.Context) error {
	if err := mw.BaseMiddleware.Handle(ctx); err != nil {
		bErr, ok := err.(*data.BaseError)
		if !ok {
			return err
		}
		switch bErr.Code {
		case data.ErrCodeValidate, data.ErrCodeFormat:
			ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusBadRequest})
		case data.ErrCodeInternal:
			ctx.SetOutput(&data.SpecCodeResponse{StatusCode: http.StatusInternalServerError})
		}
		return err
	}
	return nil
}
