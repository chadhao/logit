package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/model"
	"github.com/labstack/echo/v4"
)

func UserEntry(c echo.Context) error {
	// Check user existance
	// Create user if not existed
	// Login user if existed
	// Return token or error
	user := model.User{}

	if err := c.Bind(&user); err != nil {
		return err
	}

	if !user.Exists() && user.ValidForRegister() {
		if err := user.Create(); err != nil {
			return err
		}
	} else {
		if err := user.Login(); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func CreateDriver(c echo.Context) error {
	return nil
}

func CreateTransportOperator(c echo.Context) error {
	return nil
}
