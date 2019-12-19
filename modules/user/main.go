package user

import (
	"github.com/labstack/echo/v4"
	"github.com/chadhao/logit/config"
)

func InitMod(e *echo.Echo) {
	config.Config["test"]
	// load config
	// add routes
	// other initialization code
}
