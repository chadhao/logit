package middleware

import (
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/middleware/jwt"
	"github.com/chadhao/logit/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoadBeforeRouter(e *echo.Echo, con config.Config, r router.Router) error {
	// Routes and Config insertion
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", con)
			c.Set("router", r)
			return next(c)
		}
	})

	e.Pre(middleware.RemoveTrailingSlash())

	return nil
}

func LoadAfterRouter(e *echo.Echo, c config.Config) error {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	jwtAccessSigningKey, _ := c.Get("system.jwt.access.key")
	e.Use(jwt.JWTWithConfig(jwt.JWTConfig{
		Skipper: func(e echo.Context) bool {
			r := e.Get("router").(router.Router)
			route, err := r.Match(e.Request().Method, e.Path())
			if err != nil {
				return true
			}
			return len(route.Roles) == 0
		},
		SigningKey: []byte(jwtAccessSigningKey),
	}))

	return nil
}
