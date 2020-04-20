package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getSuscription 获取
func getSuscription(c echo.Context) error {

	req := new(reqSuscription)
	if err := c.Bind(req); err != nil {
		return err
	}

	// admin可以获取其它用户suscription, 个人只能获取自己的
	if !utils.IsOrigin(c, utils.ADMIN) {
		uid, _ := c.Get("user").(primitive.ObjectID)
		if uid.Hex() != req.DriverID {
			return errors.New("not authorized")
		}
	}

	s, err := req.getSuscription()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, s)
}
