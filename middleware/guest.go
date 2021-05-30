package middleware

import (
	"github.com/zengineDev/dojo"
	"net/http"
)

type (
	GuestConfig struct {
		Skipper      Skipper
		BeforeFunc   BeforeFunc
		RedirectPath string
	}
)

var (
	DefaultGuestConfig = GuestConfig{
		Skipper:      DefaultSkipper,
		RedirectPath: "/",
	}
)

func Guest() dojo.MiddlewareFunc {
	config := DefaultGuestConfig
	return GuestWithConfig(config)
}

func GuestWithConfig(config GuestConfig) dojo.MiddlewareFunc {

	if config.Skipper == nil {
		config.Skipper = DefaultAuthenticationConfig.Skipper
	}

	return func(next dojo.Handler) dojo.Handler {
		return func(context dojo.Context) error {
			user := context.Dojo().Auth.GetAuthUser(context)
			if !user.IsGuest() {
				http.Redirect(context.Response(), context.Request(), config.RedirectPath, http.StatusFound)
			}
			return next(context)
		}
	}
}
