package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver/backofficehandler"
	"gameapp/delivery/httpserver/matchinghandler"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config   config.Config
	uHandler userhttpserverhandler.Handler
	bHandler backofficehandler.Handler
	mHandler matchinghandler.Handler
	Router   *echo.Echo
}

func New(config config.Config, uHandler userhttpserverhandler.Handler, bHandler backofficehandler.Handler, mHandler matchinghandler.Handler) Server {
	return Server{
		config:   config,
		uHandler: uHandler,
		bHandler: bHandler,
		mHandler: mHandler,
		Router:   echo.New(),
	}
}

func (s Server) Serve() {

	// Middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	s.uHandler.SetUserRoutes(s.Router)
	s.bHandler.SetRoutes(s.Router)
	s.mHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	err := s.Router.Start(address)
	if err != nil {
		fmt.Println("err in start server", err)
	}
	//s.Router.Logger.Fatal(s.Router.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
