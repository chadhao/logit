package main

import (
	"github.com/chadhao/logit/config"
	"github.com/labstack/echo/v4"
)

var e = echo.New()
var c = config.New()

func main() {
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
