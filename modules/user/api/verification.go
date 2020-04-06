package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/labstack/echo/v4"
)

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
