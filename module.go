package main

import (
	"github.com/chadhao/logit/modules/user"
)

func loadModules() error {
	if err := user.InitModule(e, c); err != nil {
		return err
	}
	return nil
}
