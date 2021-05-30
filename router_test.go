package dojo

import (
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
)

func TestRouter_RouteGroup(t *testing.T) {
	// Add an dojo
	app := New(DefaultConfiguration{})
	r := NewRouter(app)

	r.RouteGroup("/health", func(router *Router) {
		router.Get("/status", func(ctx Context) error {
			data := make(map[string]string)
			data["m"] = "Hallo"
			return ctx.JSON(200, data)
		})
	})

	apitest.New().
		Handler(r.GetMux()).
		Get("/health/status").
		Expect(t).
		Status(http.StatusOK).
		End()

}

func TestRouter_RegisterMiddleware(t *testing.T) {
	// Add an dojo
	app := New(DefaultConfiguration{})
	r := NewRouter(app)

	// register a new middleware
	app.MiddlewareRegistry.Register("auth", func(next Handler) Handler {
		return func(context Context) error {
			data := make(map[string]string)
			data["m"] = "Middleware"
			return context.JSON(200, data)
		}
	})

	// Use the middleware on this router
	r.Use("auth")

	r.Get("/test", func(ctx Context) error {
		data := make(map[string]string)
		data["m"] = "Hallo"
		return ctx.JSON(200, data)
	})

	apitest.New().
		Handler(r.GetMux()).
		Get("/test").
		Expect(t).
		Assert(jsonpath.Equal(`$.data.m`, "Middleware")).
		Status(http.StatusOK).
		End()
}
