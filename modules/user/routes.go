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
		Handler: api.UserRegister,
	})
	r.Add(&router.Route{
		Path:    "/user",
		Method:  http.MethodGet,
		Handler: api.GetUserInfo,
	})
	r.Add(&router.Route{
		Path:    "/user",
		Method:  http.MethodPut,
		Handler: api.UserUpdate,
	})
	r.Add(&router.Route{
		Path:    "/user/driver",
		Method:  http.MethodPost,
		Handler: api.DriverRegister,
	})
	r.Add(&router.Route{
		Path:    "/user/code",
		Method:  http.MethodPost,
		Handler: api.GetVerification,
	})
	r.Add(&router.Route{
		Path:    "/email/verification",
		Method:  http.MethodPost,
		Handler: api.EmailVerify,
	})
	r.Add(&router.Route{
		Path:    "/user/forgot",
		Method:  http.MethodPost,
		Handler: api.ForgetPassword,
	})
	r.Add(&router.Route{
		Path:    "/user/vehicle",
		Method:  http.MethodPost,
		Handler: api.VehicleCreate,
	})
	r.Add(&router.Route{
		Path:    "/user/vehicle",
		Method:  http.MethodDelete,
		Handler: api.VehicleDelete,
	})
	r.Add(&router.Route{
		Path:    "/user/vehicles",
		Method:  http.MethodGet,
		Handler: api.GetVehicles,
	})
}
