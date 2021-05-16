package dojo

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type ViewAdditionalData map[string]interface{}

type ViewData struct {
	Assets []Asset
	User   Authenticable
	Data   map[string]interface{}
}

// TODO move this in the app so it is
var functions = template.FuncMap{}

func (app Application) View(ctx Context, viewName string, data ViewAdditionalData) error {
	name := filepath.Base(fmt.Sprintf("%s/%s.gohtml", app.Configuration.View.Path, viewName))
	ts, err := template.New(name).Funcs(functions).ParseFiles(fmt.Sprintf("%s/%s.gohtml", app.Configuration.View.Path, viewName))
	if err != nil {
		return err
	}

	// Load the other templates
	ts, err = ts.ParseGlob(filepath.Join(fmt.Sprintf("%s/layouts/*.gohtml", app.Configuration.View.Path)))
	if err != nil {
		return err
	}

	// TODO merge all data from the context to the view
	user := app.Auth.GetAuthUser(ctx)
	viewData := ViewData{
		Assets: app.Assets(),
		User:   &user,
		Data:   data,
	}

	err = ts.Execute(ctx.Response(), viewData)
	if err != nil {
		return err
	}

	return nil
}
