package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver/backofficehandler"
	"gameapp/delivery/httpserver/matchinghandler"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"gameapp/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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
	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogError:         true,
		LogProtocol:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}

			logger.Logger.Named("http-server").Info("request",
				zap.String("request_id", v.RequestID),
				zap.String("host", v.Host),
				zap.String("content-length", v.ContentLength),
				zap.String("protocol", v.Protocol),
				zap.String("method", v.Method),
				zap.Duration("latency", v.Latency),
				zap.String("error", errMsg),
				zap.String("remote_ip", v.RemoteIP),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	//s.Router.Use(middleware.Logger())
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
