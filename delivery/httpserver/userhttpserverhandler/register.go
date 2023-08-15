package userhttpserverhandler

import (
	"gameapp/dto"
	"gameapp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) userRegister(e echo.Context) error {
	var req dto.RegisterRequest
	bErr := e.Bind(&req)
	if bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}
	if fieldErrors, err := h.userValidator.ValidateRegisterRequest(req); err != nil {
		code, msg := httpmsg.Error(err)
		return e.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	response, rErr := h.userSvc.Register(req)
	if rErr != nil {
		code, msg := httpmsg.Error(rErr)
		return echo.NewHTTPError(code, msg)
	}
	return e.JSON(http.StatusCreated, response)
}
