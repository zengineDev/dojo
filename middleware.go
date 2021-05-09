package dojo

type MiddlewareFunc func(Handler) Handler

func (app *Application) Use(mw ...MiddlewareFunc) {
	app.Middleware.Use(mw...)
}

type MiddlewareStack struct {
	stack []MiddlewareFunc
	skips map[string]bool
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

		for _, mw := range tstack {
			h = mw(h)
		}
		return h
	}
	return h
}

func newMiddlewareStack(mws ...MiddlewareFunc) *MiddlewareStack {
	return &MiddlewareStack{
		stack: mws,
		skips: map[string]bool{},
	}
}
