package user

import (
	"github.com/labstack/echo/v4"
)

type Module struct {
}

func (m Module) InitModule(e *echo.Echo, c map[string]string) {
	// load config
	// add routes
	// other initialization code
}

func New() *Module {
	return &Module{}
}
