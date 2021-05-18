package middleware

import (
	"github.com/zengineDev/dojo"
	"net/http"
)

type (
	HttpErrorConfig struct {
		Reporter ErrorReporter
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

				application.Logger.Error(err)
				// TODO handle the errors here
				http.Error(context.Response(), "Internal server error", 500)

				// Give the error to the reporter
				reporterErr := config.Reporter(context, err)
				if reporterErr != nil {
					return reporterErr
				}

				return err
			}
			return nil
		}
	}
}
