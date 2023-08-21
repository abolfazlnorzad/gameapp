package userhttpserverhandler

import (
	"gameapp/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/users")

	userGroup.GET("/profile", h.userProfile, middleware.Auth(h.authSvc, h.authCfg), middleware.UpsertPresence(h.presenceSvc))
	userGroup.POST("/login", h.userLogin)
	userGroup.POST("/register", h.userRegister)
}
