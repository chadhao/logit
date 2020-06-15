package user

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/api"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/router"
)

func loadRoutes(r router.Router) {
	// User
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
		Roles:   []int{constant.ROLE_USER_DEFAULT},
	})
	r.Add(&router.Route{
		Path:    "/user",
		Method:  http.MethodPut,
		Handler: api.UserUpdate,
		Roles:   []int{constant.ROLE_USER_DEFAULT},
	})
	r.Add(&router.Route{
		Path:    "/user/pin",
		Method:  http.MethodPost,
		Handler: api.DriverPinCheck,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/user/code",
		Method:  http.MethodPost,
		Handler: api.GetVerification,
	})
	r.Add(&router.Route{
		Path:    "/user/code/check",
		Method:  http.MethodPost,
		Handler: api.CheckVerificationCode,
	})
	r.Add(&router.Route{
		Path:    "/email/verification",
		Method:  http.MethodGet,
		Handler: api.EmailVerify,
	})
	r.Add(&router.Route{
		Path:    "/user/forgot",
		Method:  http.MethodPost,
		Handler: api.ForgetPassword,
	})
	r.Add(&router.Route{
		Path:    "/user/users",
		Method:  http.MethodPost,
		Handler: api.UserQuery,
		Roles:   []int{constant.ROLE_SUPER, constant.ROLE_ADMIN, constant.ROLE_TO_SUPER, constant.ROLE_TO_ADMIN},
	})

	// Driver
	r.Add(&router.Route{
		Path:    "/user/driver",
		Method:  http.MethodPost,
		Handler: api.DriverRegister,
		Roles:   []int{constant.ROLE_USER_DEFAULT},
	})

	// Transport Operator
	r.Add(&router.Route{
		Path:    "/user/transportoperator",
		Method:  http.MethodPost,
		Handler: api.TransportOperatorRegister,
		Roles:   []int{constant.ROLE_USER_DEFAULT},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperators",
		Method:  http.MethodGet,
		Handler: api.GetTransportOperators,
		Roles:   []int{constant.ROLE_DRIVER, constant.ROLE_SUPER, constant.ROLE_ADMIN},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperator/drivers",
		Method:  http.MethodGet,
		Handler: api.GetDriversByTransportOperator,
		Roles:   []int{constant.ROLE_TO_SUPER, constant.ROLE_TO_ADMIN, constant.ROLE_SUPER, constant.ROLE_ADMIN},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperator/apply",
		Method:  http.MethodPost,
		Handler: api.TransportOperatorApply,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperator",
		Method:  http.MethodPut,
		Handler: api.TransportOperatorUpdate,
		Roles:   []int{constant.ROLE_TO_SUPER},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperator/identity",
		Method:  http.MethodPost,
		Handler: api.TransportOperatorAddIdentity,
		Roles:   []int{constant.ROLE_TO_SUPER},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperator/identity",
		Method:  http.MethodDelete,
		Handler: api.TransportOperatorRemoveIdentity,
		Roles:   []int{constant.ROLE_TO_SUPER},
	})
	r.Add(&router.Route{
		Path:    "/user/transportoperator/verify",
		Method:  http.MethodPost,
		Handler: api.TransportOperatorVerify,
		Roles:   []int{constant.ROLE_SUPER, constant.ROLE_ADMIN},
	})

	// Vehicle
	r.Add(&router.Route{
		Path:    "/user/vehicle",
		Method:  http.MethodPost,
		Handler: api.VehicleCreate,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/user/vehicle",
		Method:  http.MethodDelete,
		Handler: api.VehicleDelete,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/user/vehicle",
		Method:  http.MethodGet,
		Handler: api.VehicleGet,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/user/vehicles",
		Method:  http.MethodGet,
		Handler: api.GetVehicles,
		Roles:   []int{constant.ROLE_DRIVER},
	})
}
