package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddDrivingLoc 添加一条行驶信息
func AddDrivingLoc(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	reqAdd := new(reqAddDrivingLoc)
	if err := c.Bind(reqAdd); err != nil {
		return err
	}

	drivingLoc, err := reqAdd.constructToDrivingLoc(userID)
	if err != nil {
		return err
	}
	if err = drivingLoc.Save(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, drivingLoc)
}

// GetDrivingLocs 获取行驶信息
func GetDrivingLocs(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	req := new(reqDrivingLocs)
	if err := c.Bind(req); err != nil {
		return err
	}

	drivingLocs, err := req.getDrivingLocs(userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, drivingLocs)
}
