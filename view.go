package dojo

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/mux"
	"html/template"
	"path/filepath"
)

type ViewAdditionalData map[string]interface{}

type ViewData struct {
	Assets []Asset
	User   Authenticable
	Data   map[string]interface{}
}

func csrfValue(ctx Context) func() string {
	return func() string {
		// TODO make the key readable from a configuration
		return fmt.Sprintf("%s", ctx.Value("csrf"))
	}
}

func activeRoute(ctx Context) func(route string) bool {
	return func(route string) bool {
		return mux.CurrentRoute(ctx.Request()).GetName() == route
	}
}

func route(dojo *Dojo) func(name string, args ...string) string {
	muxRouter := dojo.Route.GetMux()
	return func(name string, args ...string) string {
		url, err := muxRouter.Get(name).URL(args...)
		if err != nil {
			dojo.Logger.Error(err)
			return ""
		}
		return url.String()
	}
}

func (ctx *DefaultContext) View(viewName string, data ViewAdditionalData) error {
	d := ctx.dojo
	var functions = sprig.FuncMap()
	functions["csrf"] = csrfValue(ctx)
	functions["activeRoute"] = activeRoute(ctx)
	functions["route"] = route(d)

	name := filepath.Base(fmt.Sprintf("%s/%s.gohtml", d.Configuration.View.Path, viewName))
	ts, err := template.New(name).Funcs(functions).ParseFiles(fmt.Sprintf("%s/%s.gohtml", d.Configuration.View.Path, viewName))
	if err != nil {
		return err
	}

	// Load the other templates
	ts, err = ts.ParseGlob(filepath.Join(fmt.Sprintf("%s/layouts/*.gohtml", d.Configuration.View.Path)))
	if err != nil {
		return err
	}

	ts, err = ts.ParseGlob(filepath.Join(fmt.Sprintf("%s/components/*.gohtml", d.Configuration.View.Path)))
	if err != nil {
		return err
	}

	user := d.Auth.GetAuthUser(ctx)
	viewData := ViewData{
		Assets: d.Assets(),
		User:   &user,
		Data:   data,
	}

	err = ts.Execute(ctx.Response(), viewData)
	if err != nil {
		return err
	}

	return nil
}
