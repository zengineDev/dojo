package dojo

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"sync"
)

type Context interface {
	context.Context
	Response() http.ResponseWriter
	Request() *http.Request
	Session() *Session
	Cookies() *Cookies
	Params() ParamValues
	Param(string) string
	Bind(interface{}) error
	Data() map[string]interface{}
}

type ParamValues interface {
	Get(string) string
}

func (app Application) NewContext(rc RouteConfig, w http.ResponseWriter, r *http.Request) Context {
	params := url.Values{}
	vars := mux.Vars(r)
	for k, v := range vars {
		params.Add(k, v)
	}

	// Parse URL Query String Params
	// For POST, PUT, and PATCH requests, it also parse the request body as a form.
	// Request body parameters take precedence over URL query string values in params
	if err := r.ParseForm(); err == nil {
		for k, v := range r.Form {
			for _, vv := range v {
				params.Add(k, vv)
			}
		}
	}

	session := app.getSession(r, w)

	data := &sync.Map{}

	data.Store("app", app)
	// data.Store("env", a.Env)
	// data.Store("routes", app.Routes())
	data.Store("current_route", rc)
	data.Store("current_path", r.URL.Path)
	// data.Store("contentType", ct)
	data.Store("method", r.Method)

	return &DefaultContext{
		Context: r.Context(),
		// contentType: ct,
		session:  session,
		response: w,
		request:  r,
		params:   params,
		flash:    newFlash(session),
		data:     data,
	}
}
