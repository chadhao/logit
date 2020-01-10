package user

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/api"
	"github.com/chadhao/logit/router"
)

func loadRoutes(r router.Router) {
	r.Add(&router.Route{
		Path:    "/user",
		Method:  http.MethodPost,
		Handler: api.UserEntry,
		Roles:   nil, // optional, here for impression only
	})
}
