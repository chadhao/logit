package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/suscription/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoadRoutes 路由添加
func LoadRoutes(e *echo.Echo) {
	e.GET("/suscription", getSuscription)
}

// getSuscription 获取
func getSuscription(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	s, err := model.GetSuscription(userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, s)
}
