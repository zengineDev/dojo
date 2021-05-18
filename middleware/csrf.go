package middleware

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/zengineDev/dojo"
	"github.com/zengineDev/dojo/helpers"
	"net/http"
	"strings"
)

type (
	CSRFConfig struct {
		Skipper Skipper

		TokenLength uint8 `yaml:"token_length"`

		TokenLookup string `yaml:"token_lookup"`

		ContextKey string `yaml:"context_key"`

		CookieName string `yaml:"cookie_name"`

		CookieMaxAge int `yaml:"cookie_max_age"`
	}

	csrfTokenExtractor func(dojo.Context) (string, error)
)

var (
	DefaultCSRFConfig = CSRFConfig{
		Skipper:      DefaultSkipper,
		TokenLength:  32,
		TokenLookup:  "header:" + dojo.HeaderXCSRFToken,
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}
)

func CSRF() dojo.MiddlewareFunc {
	config := DefaultCSRFConfig
	return CSRFWithConfig(config)
}

func CSRFWithConfig(config CSRFConfig) dojo.MiddlewareFunc {

	cookieStore := sessions.NewCookieStore([]byte("secret"))
	cookieStore.Options.HttpOnly = true

	// Initialize
	parts := strings.Split(config.TokenLookup, ":")
	extractor := csrfTokenFromHeader(parts[1])
	switch parts[0] {
	case "form":
		extractor = csrfTokenFromForm(parts[1])
	case "query":
		extractor = csrfTokenFromQuery(parts[1])
	}

	return func(next dojo.Handler) dojo.Handler {
		return func(context dojo.Context, application *dojo.Application) error {

			session, err := cookieStore.Get(context.Request(), config.CookieName)
			token := ""

			// Generate token
			if err != nil {
				token = helpers.RandomString(int(config.TokenLength))
			} else {
				// Reuse token
				token = fmt.Sprintf("%s", session.Values["value"])
			}

			switch context.Request().Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			default:
				clientToken, err := extractor(context)
				if err != nil {
					return err
				}
				if !validateCSRFToken(token, clientToken) {
					return dojo.NewHTTPError(http.StatusUnauthorized, "invalid csrf token")
				}
			}

			// Set CSRF in cookie
			session.Values["value"] = token
			err = session.Save(context.Request(), context.Response())
			if err != nil {
				return dojo.NewHTTPError(http.StatusBadRequest, "session save error")
			}
			context.Set(config.ContextKey, token)

			context.Response().Header().Set(dojo.HeaderVary, dojo.HeaderCookie)

			return next(context, application)
		}
	}
}

func csrfTokenFromHeader(header string) csrfTokenExtractor {
	return func(c dojo.Context) (string, error) {
		return c.Request().Header.Get(header), nil
	}
}

func csrfTokenFromForm(param string) csrfTokenExtractor {
	return func(c dojo.Context) (string, error) {
		token := c.Request().FormValue(param)
		if token == "" {
			return "", errors.New("missing csrf token in the form parameter")
		}
		return token, nil
	}
}

func csrfTokenFromQuery(param string) csrfTokenExtractor {
	return func(c dojo.Context) (string, error) {
		token := c.Params().Get(param)
		if token == "" {
			return "", errors.New("missing csrf token in the query string")
		}
		return token, nil
	}
}

func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}
