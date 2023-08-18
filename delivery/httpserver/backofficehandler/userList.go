package backofficehandler

import "github.com/labstack/echo/v4"

func (h Handler) UserList(e echo.Context) error {
	return e.JSON(200, echo.Map{
		"message": "user list received",
	})
}
