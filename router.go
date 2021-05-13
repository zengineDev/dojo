package dojo

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	middlewares []string
	router      *mux.Router
	app         *Application
}

func NewRouter(app *Application) *Router {
	return &Router{router: mux.NewRouter(), app: app}
}

func (r *Router) GetMux() *mux.Router {
	return r.router
}

// Use a registered middleware on that router
func (r *Router) Use(name string) {
	r.middlewares = append(r.middlewares, name)
}

func (r *Router) UseStack(name string) {
	stack := r.app.MiddlewareRegistry.stacks[name]
	r.middlewares = append(r.middlewares, stack...)
}

func (r *Router) Get(path string, handler Handler) {
	r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) RouteGroup(prefix string, cb func(router *Router)) {
	subRouter := r.router.PathPrefix(prefix).Subrouter()

	cb(&Router{
		router: subRouter,
		app:    r.app,
	})
}

func (r *Router) addRoute(method string, url string, h Handler) {

	// get all registered middlewares
	mws := MiddlewareStack{}
	app := r.app

	for _, mName := range r.middlewares {
		mw, err := app.MiddlewareRegistry.findMiddleware(mName)
		if err != nil {
			continue
		}
		mws.Use(mw.Handler)
	}

	config := &RouteConfig{
		Method: method,
		Path:   url,
		// HandlerName: hs,
		Handler:     h,
		App:         r.app,
		Aliases:     []string{},
		Middlewares: mws,
	}

	config.MuxRoute = r.router.Handle(url, config).Methods(method)
}

type RouteConfig struct {
	Method       string          `json:"method"`
	Path         string          `json:"path"`
	HandlerName  string          `json:"handler"`
	ResourceName string          `json:"resourceName,omitempty"`
	PathName     string          `json:"pathName"`
	Aliases      []string        `json:"aliases"`
	MuxRoute     *mux.Route      `json:"-"`
	Handler      Handler         `json:"-"`
	App          *Application    `json:"-"`
	Middlewares  MiddlewareStack `json:"-"`
}

func (r RouteConfig) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// in the route config are the middleware stack
	// these are handler we want to call before we call the route handler

	app := r.App

	c := app.NewContext(r, res, req)

	// we have now
	err := r.Middlewares.handler(r)(c, app)

	if err != nil {
		status := http.StatusInternalServerError
		//if he, ok := err.(HTTPError); ok {
		//	status = he.Status
		//}
		// things have really hit the fan if we're here!!
		//a.Logger.Error(err)
		c.Response().WriteHeader(status)
		_, err = c.Response().Write([]byte(err.Error()))
		if err != nil {
			app.Logger.Error(err)
		}
	}
}
