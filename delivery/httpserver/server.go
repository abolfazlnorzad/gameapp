package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config   config.Config
	uHandler userhttpserverhandler.Handler
}

func New(config config.Config, uHandler userhttpserverhandler.Handler) Server {
	return Server{
		config:   config,
		uHandler: uHandler,
	}
}

func (s Server) Serve() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s.uHandler.SetUserRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
