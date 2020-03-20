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
	"go.mongodb.org/mongo-driver/bson"
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
		tos    = []*model.TransportOperator{}
	)

	if err := user.Find(); err != nil {
		return err
	}

	roles := utils.RolesAssert(user.RoleIDs)
	if roles.Is(constant.ROLE_DRIVER) {
		driver.ID = uid
		driver.Find()
	}
	if roles.Is(constant.ROLE_TO_SUPER) {
		to := &model.TransportOperator{ID: uid}
		to.Find()
		tos = append(tos, to)
	}

	resp := response.UserInfoResponse{}
	resp.Format(user, driver, tos)

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

func CheckVerificationCode(c echo.Context) error {
	vr := struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}{}
	if err := c.Bind(&vr); err != nil {
		return err
	}

	red := model.Redis{Key: vr.Phone}
	code, err := red.Get()
	if err != nil {
		return err
	}
	if vr.Code != code {
		return errors.New("verification code does not match")
	}
	return c.JSON(http.StatusOK, "ok")
}

func EmailVerify(c echo.Context) error {
	er := request.EmailVerifyRequest{}
	html := "<h1>Hi there,</h1><p>Your email has been verified!</p>"

	if err := c.Bind(&er); err != nil {
		html = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return c.HTML(http.StatusBadRequest, html)
	}
	if _, err := er.Verify(); err != nil {
		html = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return c.HTML(http.StatusBadRequest, html)
	}

	return c.HTML(http.StatusOK, html)
}

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

func TransportOperatorRegister(c echo.Context) error {
	tr := request.TransportOperatorRegRequest{}

	if err := c.Bind(&tr); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	user := &model.User{ID: uid}
	if err := user.Find(); err != nil {
		return errors.New("cannot find user")
	}

	// Assign transport operator super identity
	tr.ID = uid
	if _, err := tr.Reg(); err != nil {
		return err
	}

	if !utils.RolesAssert(user.RoleIDs).Is(constant.ROLE_TO_SUPER) {
		// Update user role
		user.RoleIDs = append(user.RoleIDs, constant.ROLE_TO_SUPER)
		if err := user.Update(); err != nil {
			return err
		}
	}

	// Issue token
	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func TransportOperatorApply(c echo.Context) error {
	r := struct {
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID" query:"transportOperatorID"`
	}{}

	if err := c.Bind(&r); err != nil {
		return err
	}
	to := &model.TransportOperator{
		ID: r.TransportOperatorID,
	}
	if err := to.Find(); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	if err := to.AddDriver(uid); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "ok")
}

func GetTransportOperators(c echo.Context) error {
	r := struct {
		DriverID primitive.ObjectID `json:"driverID" query:"driverID"`
		Name     string             `json:"name" query:"name"`
	}{}

	if err := c.Bind(&r); err != nil {
		return err
	}
	filter := bson.M{}
	uid, _ := c.Get("user").(primitive.ObjectID)
	if !r.DriverID.IsZero() {
		if r.DriverID != uid {
			return errors.New("no authorization")
		}
		filter["driverIDs"] = r.DriverID
	} else {
		filter["name"] = bson.M{"$regex": "(?i)" + r.Name}
	}
	tos, err := model.FindTransportOperatorsByDriverID(filter)
	if err != nil {
		return err
	}

	tosResp := response.TransportOperatorInfoFormat(tos)
	return c.JSON(http.StatusOK, tosResp)
}

func GetVerification(c echo.Context) error {
	vr := request.VerificationRequest{}

	if err := c.Bind(&vr); err != nil {
		return err
	}

	if err := vr.Send(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "ok")
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

func VehicleCreate(c echo.Context) error {
	vr := request.VehicleCreateRequest{}

	if err := c.Bind(&vr); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_DRIVER) {
		return errors.New("is not driver")
	}

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
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	vehicle := &model.Vehicle{
		ID: vr.ID,
	}
	if err := vehicle.Find(); err != nil {
		return err
	}
	if vehicle.DriverID != uid {
		return errors.New("no authorization")
	}

	if err := vehicle.Delete(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "deleted")
}

func GetVehicles(c echo.Context) error {

	vr := struct {
		DriverID primitive.ObjectID `json:"driverID" query:"driverID"`
	}{}
	if err := c.Bind(&vr); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	if uid != vr.DriverID {
		return errors.New("no authorization")
	}

	vehicle := &model.Vehicle{
		DriverID: vr.DriverID,
	}
	vehicles, err := vehicle.FindByDriverID()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, vehicles)
}

// func UserEntry(c echo.Context) error {
// 	// Check user existance
// 	// Create user if not existed
// 	// Login user if existed
// 	// Return token or error
// 	user := model.User{}

// 	if err := c.Bind(&user); err != nil {
// 		return err
// 	}

// 	if !user.Exists() && user.ValidForRegister() {
// 		if err := user.Create(); err != nil {
// 			return err
// 		}
// 	} else {
// 		if err := user.PasswordLogin(); err != nil {
// 			return err
// 		}
// 	}

// 	token, err := user.IssueToken(c.Get("config").(config.Config))
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusCreated, token)
// }
