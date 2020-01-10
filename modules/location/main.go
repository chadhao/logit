package location

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/location/api"
	"github.com/chadhao/logit/modules/location/model"
	"github.com/chadhao/logit/router"
)

// InitModule 模块初始化
func InitModule(r router.Router, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("location")); err != nil {
		return err
	}
	api.LoadRoutes(r)
	return nil
}

// ShutdownModule 模块结束
func ShutdownModule() {
	model.Close()
}
