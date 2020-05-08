package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/router"
)

// LoadRoutes 路由添加
func LoadRoutes(r router.Router) {
	r.Add(&router.Route{
		Path:    "/logs",
		Method:  http.MethodPost,
		Handler: queryLogs,
		Roles:   []int{constant.ROLE_SUPER, constant.ROLE_ADMIN},
	})
}
