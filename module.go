package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/location"
	"github.com/chadhao/logit/modules/record"
	"github.com/chadhao/logit/modules/suscription"
	"github.com/chadhao/logit/modules/user"
	"github.com/chadhao/logit/router"
)

type moduleInit func(router.Router, config.Config) error
type moduleShutdown func()

var modulesToBeLoaded = []moduleInit{
	// message.InitModule,
	user.InitModule,

	record.InitModule,
	location.InitModule,
	suscription.InitModule,
	// ABOVE TWO MODULES NEED TO BE REFACTORED FOR THE NEW ROUTER
}

var modulesToBeShutdown = []moduleShutdown{
	user.ShutdownModule,
	record.ShutdownModule,
	location.ShutdownModule,
	suscription.ShutdownModule,
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
