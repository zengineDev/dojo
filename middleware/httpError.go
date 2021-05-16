package middleware

import (
	"github.com/zengineDev/dojo"
	"net/http"
)

type (
	HttpErrorConfig struct {
		Reporter ErrorReporter
	}

	ErrorResponder interface {
		RespondError(ctx dojo.Context, application *dojo.Application) bool
	}
)

var (
	HttpErrorDefaultConfig = HttpErrorConfig{
		Reporter: DefaultErrorReporter,
	}
)

func HttpError() dojo.MiddlewareFunc {
	config := HttpErrorDefaultConfig
	return HttpErrorWithConfig(config)
}

func HttpErrorWithConfig(config HttpErrorConfig) dojo.MiddlewareFunc {
	return func(handler dojo.Handler) dojo.Handler {
		return func(context dojo.Context, application *dojo.Application) error {
			if err := handler(context, application); err != nil {
				if er, ok := err.(ErrorResponder); ok {
					if er.RespondError(context, application) {
						return nil
					}
				}
				http.Error(context.Response(), "Internal server error", 500)

				reporterErr := config.Reporter(context, err)
				if reporterErr != nil {
					return reporterErr
				}
			}
			return nil
		}
	}
}
