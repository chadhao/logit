package api

import (
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/service"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type transportOperatorRegRequest struct {
	LicenseNumber string  `json:"licenseNumber"`
	IsCompany     bool    `json:"isCompany"`
	Name          string  `json:"name"`
	Contact       *string `json:"contact,omitempty"`
}

// TransportOperatorRegister TO组织注册
func TransportOperatorRegister(c echo.Context) error {
	req := new(transportOperatorRegRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	resp, err := service.TransportOperatorRegister(&service.TransportOperatorRegisterInput{
		Conf:          c.Get("config").(config.Config),
		UserID:        uid,
		LicenseNumber: req.LicenseNumber,
		IsCompany:     req.IsCompany,
		Name:          req.Name,
		Contact:       req.Contact,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Token)
}

// transportOperatorApplyRequest 申请加入TO组织成为司机请求参数
type transportOperatorApplyRequest struct {
	TransportOperatorID string `json:"transportOperatorID" query:"transportOperatorID"`
}

// TransportOperatorApply 申请加入TO组织成为司机
func TransportOperatorApply(c echo.Context) error {

	req := new(transportOperatorApplyRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	toID, err := primitive.ObjectIDFromHex(req.TransportOperatorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	to, err := service.TransportOperatorFind(&service.TransportOperatorFindInput{toID})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !to.IsVerified {
		return c.JSON(http.StatusBadRequest, "transport operator not verified")
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	resp, err := service.TransportOperatorIdentityAdd(&service.TransportOperatorIdentityAddInput{
		TransportOperatorID: to.ID,
		UserID:              uid,
		Identity:            model.TO_DRIVER,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.TransportOperatorIdentity)
}

// GetTransportOperators 获取tos
func GetTransportOperators(c echo.Context) error {
	driverOrigin := true
	if utils.IsOrigin(c, utils.ADMIN) {
		driverOrigin = false
	}

	resp, err := service.TransportOperatorsFind(&service.TransportOperatorsFindInput{
		IsVerified: &driverOrigin,
		IsCompany:  &driverOrigin,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resp.Tos)
}

type transportOperatorUpdateRequest struct {
	ID            primitive.ObjectID `json:"id"`
	LicenseNumber string             `json:"licenseNumber"`
	Name          string             `json:"name"`
}

// TransportOperatorUpdate TO组织更新
func TransportOperatorUpdate(c echo.Context) error {
	req := new(transportOperatorUpdateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	resp, err := service.TransportOperatorUpdate(&service.TransportOperatorUpdateInput{
		UserID:              uid,
		TransportOperatorID: req.ID,
		LicenseNumber:       req.LicenseNumber,
		Name:                req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resp.TransportOperator)
}

type transportOperatorAssignIdentityRequest struct {
	TransportOperatorID primitive.ObjectID `json:"transportOperatorID"`
	UserID              primitive.ObjectID `json:"userID"`
	Identity            model.TOIdentity   `json:"identity"`
	Contact             string             `json:"contact"`
}

// TransportOperatorAssignIdentity TO为用户添加TO_SUPER或TO_ADMIN
func TransportOperatorAssignIdentity(c echo.Context) error {

	req := new(transportOperatorAssignIdentityRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	resp, err := service.TransportOperatorAssignIdentity(&service.TransportOperatorAssignIdentityInput{
		OperatorID:          uid,
		UserID:              req.UserID,
		TransportOperatorID: req.TransportOperatorID,
		Identity:            req.Identity,
		Contact:             req.Contact,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Identity)
}

type transportOperatorRemoveIdentityRequest struct {
	TransportOperatorIdentityID string `json:"id" query:"id"`
}

// TransportOperatorRemoveIdentity TO删除角色身份
func TransportOperatorRemoveIdentity(c echo.Context) error {

	req := new(transportOperatorRemoveIdentityRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	toiID, err := primitive.ObjectIDFromHex(req.TransportOperatorIdentityID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	if err := service.TransportOperatorRemoveIdentity(&service.TransportOperatorRemoveIdentityInput{
		OperatorID:                  uid,
		TransportOperatorIdentityID: toiID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}

type transportOperatorVerifyRequest struct {
	TransportOperatorID string `json:"transportOperatorID"`
}

// TransportOperatorVerify 管理人员审批TO组织
func TransportOperatorVerify(c echo.Context) error {
	req := new(transportOperatorVerifyRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	toID, err := primitive.ObjectIDFromHex(req.TransportOperatorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := service.TransportOperatorVerify(&service.TransportOperatorVerifyInput{
		TransportOperatorID: toID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}
