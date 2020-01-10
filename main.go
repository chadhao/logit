package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/middleware"
	"github.com/chadhao/logit/router"
	"github.com/labstack/echo/v4"
)

var (
	e = echo.New()
	r = router.New()
	c = config.New()
)

func main() {
	defer shutdownModules()

	e.Debug = true
	e.HideBanner = true

	if err := c.LoadConfig(); err != nil {
		panic(err.Error())
	}

	if err := loadModules(); err != nil {
		panic(err.Error())
	}

	middleware.LoadBeforeRouter(e, r)
	middleware.LoadAfterRouter(e, c)
	r.Register(e)

	if err := e.Start("localhost:8080"); err != nil {
		panic(err.Error())
	}
}
