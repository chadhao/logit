package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/router"
)

// LoadRoutes 加载路由
func LoadRoutes(r router.Router) {
	r.Add(&router.Route{
		Path:    "/location",
		Method:  http.MethodPost,
		Handler: addDrivingLoc,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/location",
		Method:  http.MethodGet,
		Handler: getDrivingLocs,
		Roles:   []int{constant.ROLE_SUPER, constant.ROLE_ADMIN},
	})
}
