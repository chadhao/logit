package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DriverRegister(c echo.Context) error {
	dr := request.DriverRegRequest{}

	if err := c.Bind(&dr); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	roles := utils.RolesAssert(c.Get("roles"))
	if roles.Is(constant.ROLE_DRIVER) {
		return errors.New("is driver already")
	}

	user := &model.User{ID: uid}
	if err := user.Find(); err != nil {
		return errors.New("cannot find user")
	}

	// Assign driver identity
	dr.ID = uid
	if _, err := dr.Reg(); err != nil {
		return err
	}

	// Update user role and isDriver
	user.IsDriver = true
	user.RoleIDs = append(user.RoleIDs, constant.ROLE_DRIVER)
	if err := user.Update(); err != nil {
		return err
	}

	// Issue token
	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func GetDriversByTransportOperator(c echo.Context) error {
	r := struct {
		TransportOperatorID string `json:"transportOperatorID" query:"transportOperatorID"`
	}{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	transportOperatorID, err := primitive.ObjectIDFromHex(r.TransportOperatorID)
	if err != nil {
		return err
	}
	uid, _ := c.Get("user").(primitive.ObjectID)
	if !model.IsIdentity(uid, transportOperatorID, []model.TOIdentity{model.TO_SUPER, model.TO_ADMIN}) {
		return errors.New("no authorization")
	}

	toFilter := model.TransportOperatorIdentity{
		TransportOperatorID: transportOperatorID,
		Identity:            model.TO_DRIVER,
	}
	identities, err := toFilter.Filter()
	if err != nil {
		return err
	}
	driverIDs := []primitive.ObjectID{}
	for _, v := range identities {
		driverIDs = append(driverIDs, v.UserID)
	}
	drivers, err := model.GetDrivers(driverIDs)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, drivers)
}
