package middleware

import (
	"github.com/zengineDev/dojo"
	"net/http"
)

type (
	AuthenticationConfig struct {
		Skipper      Skipper
		BeforeFunc   BeforeFunc
		RedirectPath string
	}
)

var (
	DefaultAuthenticationConfig = AuthenticationConfig{
		Skipper:      DefaultSkipper,
		RedirectPath: "/login",
	}
)

func Authentication() dojo.MiddlewareFunc {
	config := DefaultAuthenticationConfig
	return AuthenticationWithConfig(config)
}

func AuthenticationWithConfig(config AuthenticationConfig) dojo.MiddlewareFunc {

	if config.Skipper == nil {
		config.Skipper = DefaultAuthenticationConfig.Skipper
	}

	return func(next dojo.Handler) dojo.Handler {
		return func(context dojo.Context) error {
			user := context.Dojo().Auth.GetAuthUser(context)
			if user.IsGuest() {
				http.Redirect(context.Response(), context.Request(), config.RedirectPath, http.StatusFound)
			}
			return next(context)
		}
	}
}
