package location

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/location/model"
	"github.com/labstack/echo/v4"
)

// InitModule 模块初始化
func InitModule(e *echo.Echo, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("location.db")); err != nil {
		return err
	}

	// add routes
	// other initialization code
	return nil
}
