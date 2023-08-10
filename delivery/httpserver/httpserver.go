package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config        config.Config
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, validator uservalidator.Validator) Server {
	return Server{
		config:        config,
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: validator,
	}
}

func (s Server) Serve() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userGroup := e.Group("/users")

	userGroup.GET("/profile", s.userProfile)
	userGroup.POST("/login", s.userLogin)
	userGroup.POST("/register", s.userRegister)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
