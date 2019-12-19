package user

import (
	"github.com/chadhao/logit/config"
	"github.com/labstack/echo/v4"
)

type Module struct {
}

func (m Module) InitModule(e *echo.Echo, c *config.Config) {
	// load config
	// add routes
	// other initialization code
}

func New() *Module {
	return &Module{}
}
