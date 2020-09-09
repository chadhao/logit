package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VehicleCreate 添加车辆信息
func VehicleCreate(c echo.Context) error {

	req := new(service.VehicleCreateInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	if req.DriverID.IsZero() {
		req.DriverID = uid
	}
	if req.DriverID != uid {
		return c.JSON(http.StatusBadRequest, "not allowed")
	}

	resp, err := service.VehicleCreate(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp.Vehicle)
}

type vehicleDeleteRequest struct {
	VehicleID primitive.ObjectID `json:"id"`
}

// VehicleDelete 删除车辆信息
func VehicleDelete(c echo.Context) error {
	req := new(vehicleDeleteRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	err := service.VehicleDelete(&service.VehicleDeleteInput{VehicleID: req.VehicleID, UserID: uid})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "deleted")
}

type vehicleGetRequest struct {
	VehicleID string `json:"id" query:"id"`
}

// VehicleGet 获取车辆信息
func VehicleGet(c echo.Context) error {

	req := new(vehicleGetRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	vehicleID, err := primitive.ObjectIDFromHex(req.VehicleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := service.VehicleGet(&service.VehicleGetInput{VehicleID: vehicleID, UserID: uid})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp.Vehicle)
}

type vehiclesGetRequest struct {
	DriverID string `json:"driverID" query:"driverID"`
}

// VehiclesGet 查询司机的车辆群信息
func VehiclesGet(c echo.Context) error {
	req := new(vehiclesGetRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	driverID, err := primitive.ObjectIDFromHex(req.DriverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	if uid != driverID {
		return c.JSON(http.StatusUnauthorized, "no authorization")
	}

	resp, err := service.VehiclesGet(&service.VehiclesGetInput{DriverID: driverID})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, resp.Vehicles)
}
