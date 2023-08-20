package userhttpserverhandler

import (
	"fmt"
	"gameapp/dto"
	"gameapp/pkg/httpmsg"
	"gameapp/service/authservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) userProfile(e echo.Context) error {
	c := getClaims(e)
	fmt.Println("claim in profile.go", c)
	response, err := h.userSvc.GetProfile(e.Request().Context(), dto.ProfileRequest{UserID: c.UserID})
	if err != nil {
		code, msg := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}
	return e.JSON(http.StatusOK, response)
}

func getClaims(e echo.Context) *authservice.Claims {
	return e.Get("user").(*authservice.Claims)
}
