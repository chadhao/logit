package user

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/api"
	"github.com/chadhao/logit/router"
)

func loadRoutes(r router.Router) {
	// r.Add(&router.Route{
	// 	Path:    "/user",
	// 	Method:  http.MethodPost,
	// 	Handler: api.UserEntry,
	// })
	r.Add(&router.Route{
		Path:    "/user/refresh",
		Method:  http.MethodPost,
		Handler: api.RefreshToken,
	})
	r.Add(&router.Route{
		Path:    "/user/existance",
		Method:  http.MethodPost,
		Handler: api.CheckExistance,
	})
	r.Add(&router.Route{
		Path:    "/user/login/password",
		Method:  http.MethodPost,
		Handler: api.PasswordLogin,
	})
	r.Add(&router.Route{
		Path:    "/user",
		Method:  http.MethodPost,
		Handler: api.DriverRegister,
	})
}
