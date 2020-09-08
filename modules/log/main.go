package message

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/log/api"
	"github.com/chadhao/logit/modules/log/model"
	"github.com/chadhao/logit/router"
)

// InitModule 模块启动
func InitModule(r router.Router, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("log")); err != nil {
		return err
	}
	api.LoadRoutes(r)
	return nil
}

// ShutdownModule 模块结束
func ShutdownModule() {
	model.Close()
}
