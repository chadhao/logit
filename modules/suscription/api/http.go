package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/modules/user/constant"
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

	roles := utils.RolesAssert(c.Get("roles"))
	switch {
	case roles.Is(constant.ROLE_ADMIN):
	case roles.Is(constant.ROLE_DRIVER):
		uid, _ := c.Get("user").(primitive.ObjectID)
		if uid != req.DriverID {
			return errors.New("not authorized")
		}
	default:
		return errors.New("not admin or driver")
	}

	s, err := req.getSuscription()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, s)
}
