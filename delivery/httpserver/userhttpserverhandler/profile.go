package userhttpserverhandler

import (
	"gameapp/dto"
	"gameapp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) userProfile(e echo.Context) error {
	authToken := e.Request().Header.Get("Authorization")
	c, err := h.authSvc.VerifyToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	response, err := h.userSvc.GetProfile(dto.ProfileRequest{UserID: c.UserID})
	if err != nil {
		code, msg := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}
	return e.JSON(http.StatusOK, response)
}
