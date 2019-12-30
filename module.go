package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/record"
	"github.com/labstack/echo/v4"
)

type moduleInit func(*echo.Echo, config.Config) error
type moduleShutdown func()

var modulesToBeLoaded = []moduleInit{
	// message.InitModule,
	// user.InitModule,
	record.InitModule,
	// location.InitModule,
}

var modulesToBeShutdown = []moduleShutdown{
	// user.ShutdownModule,
	record.ShutdownModule,
	// location.ShutdownModule,
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
