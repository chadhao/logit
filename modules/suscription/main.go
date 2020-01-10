package suscription

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/suscription/api"
	"github.com/chadhao/logit/modules/suscription/model"
	"github.com/chadhao/logit/router"
)

// InitModule 模块初始化
func InitModule(r router.Router, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("suscription.db")); err != nil {
		return err
	}

	// add routes
	api.LoadRoutes(r)
	// other initialization code
	return nil
}

// ShutdownModule 模块结束
func ShutdownModule() {
	model.Close()
}
