package main

import (
	"github.com/chadhao/logit/config"
	"github.com/labstack/echo/v4"
)

var e = echo.New()
var c = config.New()

func main() {
	c.LoadConfig()
	loadModules()
	e.Start("localhost:8080")
}
