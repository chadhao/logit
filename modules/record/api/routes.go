package api

import (
	"net/http"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/router"
)

// LoadRoutes 路由添加
func LoadRoutes(r router.Router) {
	r.Add(&router.Route{
		Path:    "/record",
		Method:  http.MethodPost,
		Handler: addRecord,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/records/sync",
		Method:  http.MethodPost,
		Handler: offlineSyncRecords,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/record/:id",
		Method:  http.MethodDelete,
		Handler: deleteLatestRecord,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/records",
		Method:  http.MethodGet,
		Handler: getRecords,
		Roles:   []int{constant.ROLE_DRIVER},
	})
	r.Add(&router.Route{
		Path:    "/record/note",
		Method:  http.MethodPost,
		Handler: addNote,
		Roles:   []int{constant.ROLE_DRIVER},
	})
}
