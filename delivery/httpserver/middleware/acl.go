package middleware

import (
	"gameapp/entity"
	"gameapp/pkg/claim"
	"gameapp/pkg/errmsg"
	"gameapp/service/aclservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Acl(service aclservice.Service, permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cl := claim.GetClaims(c)
			isAllowed, err := service.CheckAccess(cl.UserID, cl.Role, permissions...)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.SomethingWentWrong,
				})
			}
			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorMsgUserNotAllowed,
				})
			}
			return next(c)
		}
	}
}
