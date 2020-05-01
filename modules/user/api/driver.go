package api

import (
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DriverRegister 司机注册
func DriverRegister(c echo.Context) error {
	dr := request.DriverRegRequest{}

	if err := c.Bind(&dr); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	if utils.IsRole(c, constant.ROLE_DRIVER) {
		return c.JSON(http.StatusBadRequest, "is driver already")
	}

	user := &model.User{ID: uid}
	if err := user.Find(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	// Assign driver identity
	dr.ID = uid
	if _, err := dr.Reg(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Update user role and isDriver
	user.IsDriver = true
	user.RoleIDs = append(user.RoleIDs, constant.ROLE_DRIVER)
	if err := user.Update(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Issue token
	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, token)
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
