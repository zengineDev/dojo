package dojo

import (
	"github.com/gorilla/mux"
	"net/http"
)

type RouteConfig struct {
	Method       string       `json:"method"`
	Path         string       `json:"path"`
	HandlerName  string       `json:"handler"`
	ResourceName string       `json:"resourceName,omitempty"`
	PathName     string       `json:"pathName"`
	Aliases      []string     `json:"aliases"`
	MuxRoute     *mux.Route   `json:"-"`
	Handler      Handler      `json:"-"`
	App          *Application `json:"-"`
}

func (r RouteConfig) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	a := r.App

	c := a.NewContext(r, res, req)
	err := a.Middleware.handler(r)(c)

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
			a.Logger.Error(err)
		}
	}
}

func (app *Application) GET(p string, h Handler) {
	app.addRoute(http.MethodGet, p, h)
}

func (app *Application) POST(p string, h Handler) {
	app.addRoute(http.MethodPost, p, h)
}

func (app *Application) PUT(p string, h Handler) {
	app.addRoute(http.MethodPut, p, h)
}

func (app *Application) PATCH(p string, h Handler) {
	app.addRoute(http.MethodPatch, p, h)
}

func (app *Application) DELETE(p string, h Handler) {
	app.addRoute(http.MethodDelete, p, h)
}

func (app *Application) OPTIONS(p string, h Handler) {
	app.addRoute(http.MethodOptions, p, h)
}

func (app *Application) HEAD(p string, h Handler) {
	app.addRoute(http.MethodHead, p, h)
}

func (app *Application) CONNECT(p string, h Handler) {
	app.addRoute(http.MethodConnect, p, h)
}

func (app *Application) TRACE(p string, h Handler) {
	app.addRoute(http.MethodTrace, p, h)
}

func (app *Application) addRoute(method string, url string, h Handler) {

	r := &RouteConfig{
		Method: method,
		Path:   url,
		//HandlerName: hs,
		Handler: h,
		App:     app,
		Aliases: []string{},
	}

	r.MuxRoute = app.router.Handle(url, r).Methods(method)

}
