package api

import (
	"net/http"
	"time"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/service"
	"github.com/chadhao/logit/utils"
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

func GetDriversByTransportOperator(c echo.Context) error {
	r := struct {
		TransportOperatorID string `json:"transportOperatorID" query:"transportOperatorID"`
	}{}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	transportOperatorID, err := primitive.ObjectIDFromHex(r.TransportOperatorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// 若无admin权限, 则验证user是否有改TO权限
	if !utils.IsOrigin(c, "admin") {
		uid, _ := c.Get("user").(primitive.ObjectID)
		if !model.IsIdentity(uid, transportOperatorID, []model.TOIdentity{model.TO_SUPER, model.TO_ADMIN}) {
			return c.JSON(http.StatusUnauthorized, "user has no authorization")
		}
	}

	toFilter := model.TransportOperatorIdentity{
		TransportOperatorID: transportOperatorID,
		Identity:            model.TO_DRIVER,
	}
	identities, err := toFilter.Filter()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	driverIDs := []primitive.ObjectID{}
	for _, v := range identities {
		driverIDs = append(driverIDs, v.UserID)
	}
	drivers, err := model.GetDrivers(driverIDs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, drivers)
}
