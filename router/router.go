package router

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	Router interface {
		Add(*Route)
		Routes() []*Route
		Register(*echo.Echo)
		Match(string, string) (*Route, error)
	}
	router struct {
		routes []*Route
	}
	Route struct {
		Path    string
		Method  string
		Roles   []int
		Handler echo.HandlerFunc
	}
)

func (r *router) Add(route *Route) {
	route.Path = strings.ToLower(route.Path)
	r.routes = append(r.routes, route)
}

func (r *router) Routes() []*Route {
	return r.routes
}

func (r *router) Register(e *echo.Echo) {
	for _, route := range r.routes {
		e.Add(route.Method, route.Path, route.Handler)
	}
}

func (r *router) Match(m string, p string) (*Route, error) {
	for _, route := range r.routes {
		if m == route.Method && route.Path == strings.ToLower(p) {
			return route, nil
		}
	}
	return nil, errors.New("not found")
}

func New() Router {
	return &router{
		routes: make([]*Route, 0),
	}
}
