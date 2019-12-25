package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/message"
	"github.com/chadhao/logit/modules/user"
	"github.com/labstack/echo/v4"
)

type module func(*echo.Echo, config.Config) error

var modules = []module{
	message.InitModule,
	user.InitModule,
}

func loadModules() error {
	for _, m := range modules {
		if err := m(e, c); err != nil {
			return err
		}
	}

	return nil
}
