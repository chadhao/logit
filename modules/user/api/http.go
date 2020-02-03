package api

import (
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/labstack/echo/v4"
)

func CheckExistance(c echo.Context) error {
	e := request.ExistanceRequest{}

	if err := c.Bind(&e); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, e.Check())
}

func PasswordLogin(c echo.Context) error {
	l := request.LoginRequest{}

	if err := c.Bind(&l); err != nil {
		return err
	}

	user, err := l.PasswordLogin()
	if err != nil {
		return err
	}

	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func DriverRegister(c echo.Context) error {
	ur := request.UserRegistrationRequest{}
	dr := request.DriverRegistrationRequest{}

	if err := c.Bind(&ur); err != nil {
		return err
	}
	if err := c.Bind(&dr); err != nil {
		return err
	}

	// Register user
	user, err := ur.Reg()
	if err != nil {
		return err
	}

	// Register driver
	dr.Id = user.Id
	if _, err := dr.Reg(); err != nil {
		return err
	}

	// Update user role and isDriver
	user.IsDriver = true
	user.RoleIds = []int{constant.ROLE_DRIVER}
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
