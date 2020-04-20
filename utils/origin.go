package utils

import (
	"strings"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/labstack/echo/v4"
)

// ORIGIN 前端header
const ORIGIN = "X-LOGIT-ORIGIN"

// Origin 来源
type Origin string

const (
	// ADMIN 管理后台
	ADMIN Origin = "admin"
	// TRANSPORTOPERATOR TO后台
	TRANSPORTOPERATOR Origin = "transportoperator"
	// DRIVER driver前端
	DRIVER Origin = "driver"
)

// IsOrigin 从header中取出来源是否为指定来源且角色信息是否相符
func IsOrigin(c echo.Context, o Origin) bool {
	origin := c.Request().Header.Get(ORIGIN)
	if !strings.Contains(origin, string(o)) {
		return false
	}
	// 检查角色是否和header中来源相匹配
	roles := []int{}
	switch {
	case strings.Contains(origin, string(ADMIN)):
		roles = []int{constant.ROLE_SUPER, constant.ROLE_ADMIN}
	case strings.Contains(origin, string(TRANSPORTOPERATOR)):
		roles = []int{constant.ROLE_TO_SUPER, constant.ROLE_TO_ADMIN}
	case strings.Contains(origin, string(DRIVER)):
		roles = []int{constant.ROLE_TO_SUPER, constant.ROLE_DRIVER}
	}
	return AreRoles(c, roles)
}
