package user

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/router"
)

func InitModule(r router.Router, c config.Config) error {
	if err := model.New(c.LoadModuleConfig("user")); err != nil {
		return err
	}

	loadRoutes(r)

	return nil
}

func ShutdownModule() {
	model.Close()
}
