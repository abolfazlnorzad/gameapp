package middleware

import (
	"gameapp/dto"
	"gameapp/pkg/claim"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/timestamp"
	"gameapp/service/precenseservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpsertPresence(service precenseservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cl := claim.GetClaims(c)
			_, err := service.Upsert(c.Request().Context(), dto.UpsertPresenceRequest{UserID: cl.UserID, Timestamp: timestamp.Now()})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.SomethingWentWrong,
				})
			}
			return next(c)
		}
	}
}
