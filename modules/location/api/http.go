package api

import (
	"net/http"
	"time"

	"github.com/chadhao/logit/modules/location/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// reqAddDrivingLoc 添加行驶信息请求结构
type addDrivingLocRequest struct {
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// addDrivingLoc 添加一条行驶信息
func addDrivingLoc(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)

	req := new(addDrivingLocRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := service.CreateDrivingLoc(&service.CreateDrivingLocInput{
		DriverID:  uid,
		Lat:       req.Lat,
		Lng:       req.Lng,
		CreatedAt: req.CreatedAt,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

// reqDrivingLocs 行驶信息请求结构
type findDrivingLocsRequest struct {
	service.FindDrivingLocsInput
}

// findDrivingLocs 获取行驶信息
func findDrivingLocs(c echo.Context) error {
	req := new(service.FindDrivingLocsInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := service.FindDrivingLocs(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resp.DrivingLocs)
}
