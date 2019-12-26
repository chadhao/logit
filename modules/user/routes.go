package user

import (
	"github.com/chadhao/logit/modules/user/api"
	"github.com/labstack/echo/v4"
)

func loadRoutes(e *echo.Echo) {
	e.POST("/user", api.UserEntry)
}
