package dojo

import (
	"io/ioutil"
	"log"
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

	files, err := ioutil.ReadDir(app.Configuration.Assets.Path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		parts := strings.Split(file.Name(), ".")

		assets = append(assets, Asset{
			Name:      parts[0],
			Extension: FileExtension(parts[1]),
		})
	}

	return assets
}
