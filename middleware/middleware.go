package middleware

import (
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/middleware/jwt"
	"github.com/chadhao/logit/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoadBeforeRouter(e *echo.Echo, con config.Config, r router.Router) {
	// Routes and Config injection
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", con)
			c.Set("router", r)
			return next(c)
		}
	})

	e.Pre(middleware.RemoveTrailingSlash())
}

func LoadAfterRouter(e *echo.Echo, c config.Config) {
	// CORS handling
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	// JWT handling
	jwtAccessSigningKey, _ := c.Get("system.jwt.access.key")

	e.Use(jwt.JWTWithConfig(jwt.JWTConfig{
		Skipper: func(e echo.Context) bool {
			r := e.Get("router").(router.Router)
			if _, err := r.Match(e.Request().Method, e.Path()); err != nil {
				return true
			}
			return false
			// route, err := r.Match(e.Request().Method, e.Path())
			// if err != nil {
			// 	return true
			// }
			// return len(route.Roles) == 0
		},
		SigningKey: []byte(jwtAccessSigningKey),
	}))

	//Authorization
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			r := e.Get("router").(router.Router)
			route, err := r.Match(e.Request().Method, e.Path())
			if err != nil {
				return err
			}

			if len(route.Roles) == 0 {
				return next(e)
			}

			userRoles := e.Get("roles").([]int)
			if hasIntersectionInt(route.Roles, userRoles) {
				return next(e)
			}

			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	})
}

func hasIntersectionInt(a []int, b []int) bool {
	i := make(map[int]bool, 0)

	for _, v := range a {
		i[v] = true
	}

	for _, v := range b {
		if _, ok := i[v]; ok {
			return true
		}
	}

	return false
}
