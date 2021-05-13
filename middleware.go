package dojo

import (
	"fmt"
)

type MiddlewareFunc func(Handler) Handler

type MiddlewareStack struct {
	stack []MiddlewareFunc
}

func (ms *MiddlewareStack) Use(mw ...MiddlewareFunc) {
	ms.stack = append(ms.stack, mw...)
}

func (ms *MiddlewareStack) handler(rc RouteConfig) Handler {
	h := rc.Handler
	if len(ms.stack) > 0 {
		mh := func(_ Handler) Handler {
			return h
		}
		tstack := []MiddlewareFunc{mh}
		sl := len(ms.stack) - 1
		for i := sl; i >= 0; i-- {
			mw := ms.stack[i]
			tstack = append(tstack, mw)
		}

		for _, mw := range tstack {
			h = mw(h)
		}
		return h
	}
	return h
}

type Middleware struct {
	Name    string
	Handler MiddlewareFunc
}

type MiddlewareRegistry struct {
	middlewares []Middleware
	stacks      map[string][]string
}

func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{
		stacks: make(map[string][]string),
	}
}

func (registry *MiddlewareRegistry) Register(name string, fn MiddlewareFunc) {
	registry.middlewares = append(registry.middlewares, Middleware{
		Name:    name,
		Handler: fn,
	})
}

func (registry *MiddlewareRegistry) RegisterStack(name string, middlewares []string) {
	registry.stacks[name] = middlewares
}

func (registry MiddlewareRegistry) findMiddleware(name string) (Middleware, error) {
	for _, m := range registry.middlewares {
		if m.Name == name {
			return m, nil
		}
	}
	return Middleware{}, fmt.Errorf("middleware %s is not registered", name)
}
