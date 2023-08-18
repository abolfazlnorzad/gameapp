package backofficehandler

import (
	"gameapp/delivery/httpserver/middleware"
	"gameapp/entity"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	e.GET("/bc", h.UserList, middleware.Auth(h.authSvc, h.authCfg), middleware.Acl(h.aclSvc, entity.UserListPermission))
}
