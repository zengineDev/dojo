package middleware

import (
	"github.com/zengineDev/dojo"
	"time"
)

type (
	LoggingConfig struct {
		Skipper Skipper
	}
)

var (
	DefaultLoggingConfig = LoggingConfig{
		Skipper: DefaultSkipper,
	}
)

func Logging() dojo.MiddlewareFunc {
	config := DefaultLoggingConfig
	return LoggingWithConfig(config)
}

func LoggingWithConfig(config LoggingConfig) dojo.MiddlewareFunc {
	return func(next dojo.Handler) dojo.Handler {
		return func(context dojo.Context, application *dojo.Application) error {
			fields := make(map[string]interface{})

			start := time.Now()
			if err := next(context, application); err != nil {
				// TODO add better error handling
			}
			stop := time.Now()
			p := context.Request().URL.Path
			if p == "" {
				p = "/"
			}
			fields["path"] = p
			fields["method"] = context.Request().Method
			fields["uri"] = context.Request().RequestURI
			fields["host"] = context.Request().Host
			fields["env"] = application.Configuration.App.Environment
			id := context.Request().Header.Get(dojo.HeaderXRequestID)
			if id == "" {
				id = context.Response().Header().Get(dojo.HeaderXRequestID)
			}
			fields["id"] = id
			fields["remote_ip"] = context.RealIP()
			fields["start"] = start
			fields["stop"] = stop
			fields["contentType"] = context.Request().Header.Get(dojo.HeaderContentType)
			fields["userAgent"] = context.Request().UserAgent()
			fields["referer"] = context.Request().Referer()
			fields["protocol"] = context.Request().Proto
			//n := context.Response().Status()
			n := 200
			var s int
			switch {
			case n >= 500:
				s = n
			case n >= 400:
				s = n
			case n >= 300:
				s = n
			}
			fields["status"] = s

			application.Logger.WithFields(fields).Info("request_log")
			return nil
		}
	}
}
