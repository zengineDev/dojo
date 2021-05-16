package middleware

import "github.com/zengineDev/dojo"

type (
	Skipper       func(ctx dojo.Context) bool
	BeforeFunc    func(ctx dojo.Context)
	ErrorReporter func(ctx dojo.Context, err error) error
)

func DefaultSkipper(ctx dojo.Context) bool {
	return false
}

func DefaultErrorReporter(ctx dojo.Context, err error) error {
	return nil
}
