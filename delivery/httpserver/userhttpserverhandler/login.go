package userhttpserverhandler

import (
	"gameapp/dto"
	"gameapp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) userLogin(e echo.Context) error {
	var req dto.LoginRequest
	bErr := e.Bind(&req)
	if bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	if fieldErrors, err := h.userValidator.ValidateLoginRequest(req); err != nil {
		code, msg := httpmsg.Error(err)
		return e.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	response, err := h.userSvc.Login(req)
	if err != nil {
		code, msg := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}
	return e.JSON(http.StatusOK, response)
}
