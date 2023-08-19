package matchinghandler

import (
	"gameapp/dto"
	"gameapp/pkg/claim"
	"gameapp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) AddToWaitList(e echo.Context) error {
	var req dto.AddToWaitingListRequest
	if err := e.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	claims := claim.GetClaims(e)
	req.UserID = claims.UserID

	if fErr, err := h.matchingValidator.ValidateAddToWaitingListRequest(req); err != nil {
		code, msg := httpmsg.Error(err)
		return e.JSON(code, echo.Map{
			"message": msg,
			"errors":  fErr,
		})
	}

	res, err := h.matchingSvc.AddToWaitingList(req)
	if err != nil {
		code, msg := httpmsg.Error(err)
		return e.JSON(code, msg)
	}

	return e.JSON(http.StatusOK, res)

}
