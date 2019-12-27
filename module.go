package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/message"
	"github.com/chadhao/logit/modules/user"
	"github.com/labstack/echo/v4"
)

type moduleInit func(*echo.Echo, config.Config) error
type moduleShutdown func()

var modulesToBeLoaded = []moduleInit{
	message.InitModule,
	user.InitModule,
}

var modulesToBeShutdown = []moduleShutdown{
	user.ShutdownModule,
}

func loadModules() error {
	for _, m := range modulesToBeLoaded {
		if err := m(e, c); err != nil {
			return err
		}
	}

	return nil
}

func shutdownModules() {
	for _, m := range modulesToBeShutdown {
		m()
	}
}
