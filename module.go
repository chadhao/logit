package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/message"
	"github.com/chadhao/logit/modules/user"
	"github.com/chadhao/logit/router"
)

type moduleInit func(router.Router, config.Config) error
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
		if err := m(r, c); err != nil {
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
