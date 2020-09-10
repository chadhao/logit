package api

import (
	"net/http"
	"time"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// driverRegisterRequest 司机注册请求参数
type driverRegisterRequest struct {
	LicenseNumber string    `json:"licenseNumber" valid:"stringlength(5|8)`
	DateOfBirth   time.Time `json:"dateOfBirth" valid:"required"`
	Firstnames    string    `json:"firstnames" valid:"required"`
	Surname       string    `json:"surname" valid:"required"`
	Pin           string    `json:"pin" valid:"stringlength(4|4)`
}

// DriverRegister 司机注册, 司机注册时需要添加pin码
func DriverRegister(c echo.Context) error {

	req := new(driverRegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	resp, err := service.DriverRegister(&service.DriverRegisterInput{
		Conf:          c.Get("config").(config.Config),
		UserID:        uid,
		LicenseNumber: req.LicenseNumber,
		DateOfBirth:   req.DateOfBirth,
		Firstnames:    req.Firstnames,
		Surname:       req.Surname,
		Pin:           req.Pin,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Token)
}

type driverPinCheckRequest struct {
	Pin string `json:"pin"`
}

// DriverPinCheck 司机验证pin码
func DriverPinCheck(c echo.Context) error {
	req := new(driverPinCheckRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	if err := service.DriverPinCheck(&service.DriverPinCheckInput{UserID: uid, Pin: req.Pin}); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}

type getDriversByTORequest struct {
	TransportOperatorID string `json:"transportOperatorID" query:"transportOperatorID"`
}

// GetDriversByTransportOperator 获取TO组织下的司机信息
func GetDriversByTransportOperator(c echo.Context) error {
	req := new(getDriversByTORequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	toID, err := primitive.ObjectIDFromHex(req.TransportOperatorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	resp, err := service.DriversFindByTO(&service.DriversFindByTOInput{
		OperatorID:          uid,
		TransportOperatorID: toID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Drivers)
}
