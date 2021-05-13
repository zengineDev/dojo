package dojo

import (
	"os"
	"strings"
)

// get the file name and the extension to put it in the template

type Asset struct {
	Name      string
	Extension FileExtension
}

type FileExtension string

const (
	Javascript FileExtension = "js"
	Stylesheet FileExtension = "css"
)

func (app *Application) Assets() []Asset {
	var assets []Asset

	files, err := os.ReadDir(app.Configuration.Assets.Path)
	if err != nil {
		app.Logger.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		app.Logger.Debugf("assets: register file %s", file.Name())
		parts := strings.Split(file.Name(), ".")
		if len(parts) > 1 {
			assets = append(assets, Asset{
				Name:      parts[0],
				Extension: FileExtension(parts[1]),
			})
		}
	}

	return assets
}
