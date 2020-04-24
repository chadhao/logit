package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/chadhao/logit/modules/user/response"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckExistance(c echo.Context) error {
	r := request.ExistanceRequest{}

	if err := c.Bind(&r); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r.Check())
}

func RefreshToken(c echo.Context) error {
	r := request.RefreshTokenRequest{}

	if err := c.Bind(&r); err != nil {
		return err
	}

	config := c.Get("config").(config.Config)
	user, err := r.Validate(config)
	if err != nil {
		return err
	}
	if err := user.Find(); err != nil {
		return err
	}
	token, err := user.IssueToken(config)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func PasswordLogin(c echo.Context) error {
	r := request.LoginRequest{}

	if err := c.Bind(&r); err != nil {
		return err
	}

	user, err := r.PasswordLogin()
	if err != nil {
		return err
	}

	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func GetUserInfo(c echo.Context) error {
	uid, _ := c.Get("user").(primitive.ObjectID)

	var (
		user   = &model.User{ID: uid}
		driver = &model.Driver{}
	)

	if err := user.Find(); err != nil {
		return err
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

func UserRegister(c echo.Context) error {
	ur := request.UserRegRequest{}

	if err := c.Bind(&ur); err != nil {
		return err
	}

	user, err := ur.Reg()
	if err != nil {
		return err
	}

	// Issue token
	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return err
	}

	// 当用户注册后为用户发送email验证邮件
	go func(email string) {
		vr := request.VerificationRequest{
			Email: email,
		}
		vr.Send()
	}(ur.Email)
	return c.JSON(http.StatusOK, token)
}

func UserUpdate(c echo.Context) error {
	ur := request.UserUpdateRequest{}

	if err := c.Bind(&ur); err != nil {
		return err
	}
	uid, _ := c.Get("user").(primitive.ObjectID)
	user := &model.User{ID: uid}
	if err := user.Find(); err != nil {
		return err
	}

	if err := ur.Replace(user); err != nil {
		return err
	}

	if err := user.Update(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "ok")
}

func ForgetPassword(c echo.Context) error {
	vr := request.ForgetPasswordRequest{}

	if err := c.Bind(&vr); err != nil {
		return err
	}

	user := model.User{Phone: vr.Phone, Email: vr.Email}
	if err := user.Find(); err != nil {
		return err
	}

	if err := vr.Verify(); err != nil {
		return err
	}

	user.Password = vr.Password
	if err := user.Update(); err != nil {
		return err
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
		return err
	}

	if len(ur.Phone) > 0 && len(ur.Phone) < 4 {
		return errors.New("phone number need to be more specific")
	}
	if len(ur.Email) > 0 && len(ur.Email) < 4 {
		return errors.New("email need to be more specific")
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
