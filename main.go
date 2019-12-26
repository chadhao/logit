package main

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/middleware"
	"github.com/labstack/echo/v4"
)

var (
	e = echo.New()
	c = config.New()
)

func main() {
	defer shutdownModules()

	e.Debug = true
	e.HideBanner = true
	middleware.LoadBeforeRouter(e)
	middleware.LoadAfterRouter(e)

	if err := c.LoadConfig(); err != nil {
		panic(err.Error())
	}

	if err := loadModules(); err != nil {
		panic(err.Error())
	}

	if err := e.Start("localhost:8080"); err != nil {
		panic(err.Error())
	}
}
