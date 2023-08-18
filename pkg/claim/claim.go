package claim

import (
	"gameapp/service/authservice"
	"github.com/labstack/echo/v4"
)

func GetClaims(e echo.Context) *authservice.Claims {
	return e.Get("user").(*authservice.Claims)
}
