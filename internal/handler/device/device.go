package device

import (
	"giot/internal/core/model"
	"giot/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"
	"reflect"
)

type Handler struct {
}

func NewHandler() handler.RouteRegister {
	return &Handler{}
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/giot/device", wgin.Wraps(h.Create,
		wrapper.InputType(reflect.TypeOf(model.Device{}))))
}
func (h *Handler) Create(c droplet.Context) (interface{}, error) {
	input := c.Input().(*model.Device)
	var result, _ = stroage.DB.SqlMapClient("query_history", "guest").Query().List()

	return result, nil
}
