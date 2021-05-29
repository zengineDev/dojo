package dojo

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	middlewares []string
	router      *mux.Router
	dojo        *Dojo
}

func NewRouter(dojo *Dojo) *Router {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/dist"))))
	return &Router{router: r, dojo: dojo}
}

func (r *Router) GetMux() *mux.Router {
	return r.router
}

// Use a registered middleware on that router
func (r *Router) Use(name string) {
	r.middlewares = append(r.middlewares, name)
}

func (r *Router) UseStack(name string) {
	stack := r.dojo.MiddlewareRegistry.stacks[name]
	r.middlewares = append(r.middlewares, stack...)
}

func (r *Router) Get(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodGet, path, handler, middlewares...)
}

func (r *Router) GetWithName(path string, name string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addNamedRoute(http.MethodGet, path, name, handler, middlewares...)
}

func (r *Router) Post(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodPost, path, handler, middlewares...)
}

func (r *Router) PostWithName(path string, name string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addNamedRoute(http.MethodPost, path, name, handler, middlewares...)
}

func (r *Router) Put(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodPut, path, handler, middlewares...)
}

func (r *Router) PutWithName(path string, name string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addNamedRoute(http.MethodPut, path, name, handler, middlewares...)
}

func (r *Router) Patch(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodPatch, path, handler, middlewares...)
}

func (r *Router) PatchWithName(path string, name string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addNamedRoute(http.MethodPatch, path, name, handler, middlewares...)
}

func (r *Router) Options(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodOptions, path, handler, middlewares...)
}

func (r *Router) Delete(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodDelete, path, handler, middlewares...)
}

func (r *Router) DeleteWithName(path string, name string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addNamedRoute(http.MethodDelete, path, name, handler, middlewares...)
}

func (r *Router) Trace(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodTrace, path, handler, middlewares...)
}

func (r *Router) Connect(path string, handler Handler, middlewares ...MiddlewareFunc) {
	r.addRoute(http.MethodConnect, path, handler, middlewares...)
}

func (r *Router) RouteGroup(prefix string, cb func(router *Router)) {
	subRouter := r.router.PathPrefix(prefix).Subrouter()

	cb(&Router{
		router: subRouter,
		dojo:   r.dojo,
	})
}

func (r *Router) getRouteConfig(method string, url string, h Handler) RouteConfig {
	mws := MiddlewareStack{}
	app := r.dojo

	for _, mName := range r.middlewares {
		mw, err := app.MiddlewareRegistry.findMiddleware(mName)
		if err != nil {
			continue
		}
		mws.Use(mw.Handler)
	}

	return RouteConfig{
		Method: method,
		Path:   url,
		// HandlerName: hs,
		Handler:     h,
		Dojo:        r.dojo,
		Aliases:     []string{},
		Middlewares: mws,
	}
}

func (r *Router) addNamedRoute(method string, url string, name string, h Handler, middlewares ...MiddlewareFunc) {
	config := r.getRouteConfig(method, url, h)
	config.Middlewares.Use(middlewares...)
	config.MuxRoute = r.router.Handle(url, config).Methods(method).Name(name)
}

func (r *Router) addRoute(method string, url string, h Handler, middlewares ...MiddlewareFunc) {
	config := r.getRouteConfig(method, url, h)
	config.Middlewares.Use(middlewares...)
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
	Dojo         *Dojo           `json:"-"`
	Middlewares  MiddlewareStack `json:"-"`
}

func (r RouteConfig) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app := r.Dojo
	c := app.NewContext(r, res, req)
	err := r.Middlewares.handler(r)(c)
	if err != nil {
		app.HTTPErrorHandler(err, c)
	}
}

func (r Router) Redirect(ctx Context, url string) {
	http.Redirect(ctx.Response(), ctx.Request(), url, 302)
}
