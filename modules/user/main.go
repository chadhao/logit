package user

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/labstack/echo/v4"
)

func InitModule(e *echo.Echo, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("user.db")); err != nil {
		return err
	}

	// add routes
	// other initialization code

	return nil
}

func ShutdownModule() {
	model.Close()
}
