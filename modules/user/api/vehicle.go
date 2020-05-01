package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VehicleCreate(c echo.Context) error {
	vr := request.VehicleCreateRequest{}

	if err := c.Bind(&vr); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	vr.DriverID = uid
	vehicle, err := vr.Create()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, vehicle)
}

func VehicleDelete(c echo.Context) error {

	vr := struct {
		ID primitive.ObjectID `json:"id"`
	}{}
	if err := c.Bind(&vr); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	vehicle := &model.Vehicle{
		ID: vr.ID,
	}
	if err := vehicle.Find(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	if vehicle.DriverID != uid {
		return c.JSON(http.StatusUnauthorized, "no authorization")
	}

	if err := vehicle.Delete(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "deleted")
}

func VehicleGet(c echo.Context) error {

	vr := struct {
		ID string `json:"id" query:"id"`
	}{}
	if err := c.Bind(&vr); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	vid, err := primitive.ObjectIDFromHex(vr.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	vehicle := &model.Vehicle{
		ID: vid,
	}
	if err := vehicle.Find(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	if vehicle.DriverID != uid {
		return c.JSON(http.StatusUnauthorized, "no authorization")
	}

	return c.JSON(http.StatusOK, vehicle)
}

func GetVehicles(c echo.Context) error {
	vr := struct {
		DriverID string `json:"driverID" query:"driverID"`
	}{}
	if err := c.Bind(&vr); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	driverID, err := primitive.ObjectIDFromHex(vr.DriverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)
	if uid != driverID {
		return c.JSON(http.StatusUnauthorized, "no authorization")
	}
	vehicle := &model.Vehicle{
		DriverID: driverID,
	}
	vehicles, err := vehicle.FindByDriverID()
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, vehicles)
}
