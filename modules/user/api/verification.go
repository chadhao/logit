package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/service"
	"github.com/labstack/echo/v4"
)

// CheckVerificationCode 手机验证码验证
func CheckVerificationCode(c echo.Context) error {
	req := new(service.CheckVerificationCodeInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := service.CheckVerificationCode(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}

// EmailVerify 邮箱验证
func EmailVerify(c echo.Context) error {
	req := new(service.EmailVerifyInput)
	if err := c.Bind(req); err != nil {
		return c.HTML(http.StatusBadRequest, "<h1>Bad request</h1><p>"+err.Error()+"</p>")
	}

	out, err := service.EmailVerify(req)
	if err != nil {
		return c.HTML(http.StatusBadRequest, out.HTML)
	}

	return c.HTML(http.StatusOK, out.HTML)
}

// GetVerification 发送手机或邮箱验证码
func GetVerification(c echo.Context) error {

	req := new(service.SendVerificationInput)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := service.SendVerification(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}
