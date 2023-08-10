package httpserver

import (
	"fmt"
	"gameapp/pkg/httpmsg"
	"gameapp/service/userservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) userLogin(e echo.Context) error {
	var req userservice.LoginRequest
	bErr := e.Bind(&req)
	if bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	fmt.Println("req ", req)
	response, err := s.userSvc.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, response)
}

func (s Server) userRegister(e echo.Context) error {
	var req userservice.RegisterRequest
	bErr := e.Bind(&req)
	if bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	response, rErr := s.userSvc.Register(req)
	if rErr != nil {
		code, msg := httpmsg.Error(rErr)
		return echo.NewHTTPError(code, msg)
	}
	return e.JSON(http.StatusCreated, response)
}

func (s Server) userProfile(e echo.Context) error {
	authToken := e.Request().Header.Get("Authorization")
	c, err := s.authSvc.VerifyToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	response, err := s.userSvc.GetProfile(userservice.ProfileRequest{UserID: c.UserID})
	if err != nil {
		code, msg := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}
	return e.JSON(http.StatusOK, response)
}
