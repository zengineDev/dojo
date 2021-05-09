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
		c.Response().Write([]byte(err.Error()))
	}
}

func (app *Application) GET(p string, h Handler) {
	app.addRoute(http.MethodGet, p, h)
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
