package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoadRoutes 路由添加
func LoadRoutes(e *echo.Echo) {
	e.GET("/suscription", getSuscription)
}

// getSuscription 获取
func getSuscription(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
