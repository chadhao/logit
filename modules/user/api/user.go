package api

import (
	"net/http"

	"github.com/chadhao/logit/config"
	logInternals "github.com/chadhao/logit/modules/log/internals"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/chadhao/logit/modules/user/response"
	"github.com/chadhao/logit/modules/user/service"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckUserExistance 检查用户是否存在
func CheckUserExistance(c echo.Context) error {

	req := new(service.UserExistanceCheckInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp := service.UserExistanceCheck(req)
	return c.JSON(http.StatusOK, resp)
}

// refreshTokenRequest 更新token参数
type refreshTokenRequest struct {
	Token string `json:"token"`
}

// RefreshToken 更新token
func RefreshToken(c echo.Context) error {

	req := new(refreshTokenRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	config := c.Get("config").(config.Config)
	resp, err := service.RefreshToken(&service.RefreshTokenInput{Token: req.Token, Conf: config})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Token)
}

// passwordLoginRequest 密码登录参数
type passwordLoginRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	License  string `json:"license"`
	Password string `json:"password"`
}

// PasswordLogin 密码登录
func PasswordLogin(c echo.Context) error {
	req := new(passwordLoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := service.PasswordLogin(&service.PasswordLoginInput{
		Phone:    req.Phone,
		Email:    req.Email,
		License:  req.License,
		Password: req.Password,
		Conf:     c.Get("config").(config.Config),
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp.Token)
}

// UserInfoGet 获取用户信息
func UserInfoGet(c echo.Context) error {
	uid, _ := c.Get("user").(primitive.ObjectID)

	var (
		user   = &model.User{ID: uid}
		driver = &model.Driver{}
	)

	if err := user.Find(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	roles := utils.RolesAssert(user.RoleIDs)
	if roles.Is(constant.ROLE_DRIVER) {
		driver.ID = uid
		driver.Find()
	}

	tFilter := model.TransportOperatorIdentity{
		UserID: uid,
	}
	tois, _ := tFilter.Filter()

	resp := response.UserInfoResponse{}
	resp.Format(user, driver, tois)

	return c.JSON(http.StatusOK, resp)
}

// userRegisterRequest 用户注册请求参数
type userRegisterRequest struct {
	*service.UserRegisterInput
}

// UserRegister 用户注册
func UserRegister(c echo.Context) error {
	req := new(userRegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userResp, err := service.UserRegister(req.UserRegisterInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Issue token
	tokenResp, err := service.IssueToken(&service.IssueTokenInput{
		UserID:  userResp.ID,
		RoleIDs: userResp.RoleIDs,
		Conf:    c.Get("config").(config.Config),
	})
	if err != nil {
		return err
	}

	// 当用户注册后为用户发送email验证邮件
	go func(email string) {
		vr := request.VerificationRequest{
			Email: email,
		}
		vr.Send()
	}(userResp.Email)
	return c.JSON(http.StatusOK, tokenResp.Token)
}

// userUpdateRequest 用户更新请求参数
type userUpdateRequest struct {
	Password *string `json:"password,omitempty"`
	Pin      *string `json:"pin,omitempty"`
}

// UserUpdate 用户更新
func UserUpdate(c echo.Context) error {

	req := new(userUpdateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	if err := service.UserUpdate(&service.UserUpdateInput{
		UserID:   uid,
		Password: req.Password,
		Pin:      req.Pin,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// log 记录
	go func(from primitive.ObjectID, content interface{}) {
		log := &logInternals.AddLogRequest{
			Type:    "modification",
			FromFun: "UserUpdate",
			From:    &from,
			Content: content,
		}
		logInternals.AddLog(log)
	}(uid, *req)

	return c.JSON(http.StatusOK, "ok")
}

// ForgetPassword 忘记密码
func ForgetPassword(c echo.Context) error {

	req := new(service.ForgetPasswordInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := service.ForgetPassword(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}

// UserQuery 根据条件查询用户
func UserQuery(c echo.Context) error {
	ur := struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
	}{}

	if err := c.Bind(&ur); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(ur.Phone) > 0 && len(ur.Phone) < 4 {
		return c.JSON(http.StatusBadRequest, "phone number need to be more specific")
	}
	if len(ur.Email) > 0 && len(ur.Email) < 4 {
		return c.JSON(http.StatusBadRequest, "email need to be more specific")
	}

	userFilter := model.User{
		Phone: ur.Phone,
		Email: ur.Email,
	}

	users, err := userFilter.Filter()
	if err != nil {
		return err
	}

	uids := []primitive.ObjectID{}
	for _, v := range users {
		uids = append(uids, v.ID)
	}
	driversMap := make(map[primitive.ObjectID]model.Driver)
	drivers, _ := model.GetDrivers(uids)
	for _, driver := range drivers {
		driversMap[driver.ID] = driver
	}

	identitiesMap, _ := model.GetIdentitiesByUserIDs(uids)

	type resp struct {
		model.User
		Driver     model.Driver                            `json:"driver"`
		Identities []model.TransportOperatorIdentityDetail `json:"identities"`
	}
	resps := []resp{}
	for _, user := range users {
		r := resp{
			User:       user,
			Driver:     driversMap[user.ID],
			Identities: identitiesMap[user.ID],
		}
		resps = append(resps, r)
	}

	return c.JSON(http.StatusOK, resps)
}
