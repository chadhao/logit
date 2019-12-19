package main

import (
	echo "github.com/labstack/echo/v4"
)

var modules = make([]interface{}, 0)

func main() {
	e := echo.New()
	e.Start("localhost:8080")
}
