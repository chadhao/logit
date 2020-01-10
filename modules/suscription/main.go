package suscription

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/suscription/api"
	"github.com/chadhao/logit/modules/suscription/model"
	"github.com/labstack/echo/v4"
)

// InitModule 模块初始化
func InitModule(e *echo.Echo, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("suscription.db")); err != nil {
		return err
	}

	// add routes
	api.LoadRoutes(e)
	// other initialization code
	return nil
}

// ShutdownModule 模块结束
func ShutdownModule() {
	model.Close()
}
