package matchinghandler

import (
	"gameapp/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	e.POST("matching/add-to-wait-list", h.AddToWaitList, middleware.Auth(h.authSvc, h.authCfg))
}
