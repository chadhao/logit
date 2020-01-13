package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// addDrivingLoc 添加一条行驶信息
func addDrivingLoc(c echo.Context) error {

	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_DRIVER) {
		return errors.New("not driver")
	}
	userID, _ := c.Get("user").(primitive.ObjectID)

	req := new(reqAddDrivingLoc)
	if err := c.Bind(req); err != nil {
		return err
	}

	drivingLoc, err := req.constructToDrivingLoc(userID)
	if err != nil {
		return err
	}
	if err = drivingLoc.Save(); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, drivingLoc)
}

// getDrivingLocs 获取行驶信息
func getDrivingLocs(c echo.Context) error {

	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_SUPER) {
		return errors.New("not super")
	}

	req := new(reqDrivingLocs)
	if err := c.Bind(req); err != nil {
		return err
	}

	drivingLocs, err := req.getDrivingLocs()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, drivingLocs)
}
