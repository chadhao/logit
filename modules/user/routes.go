package user

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/api"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/router"
)

func loadRoutes(r router.Router) {
	r.Add(&router.Route{
		Path:    "/user",
		Method:  http.MethodPost,
		Handler: api.UserEntry,
		Roles:   []int{constant.ROLE_ADMIN}, // optional, here for impression only
	})
}
