package dojo

import (
	"fmt"
	"os"
	"strings"
)

// get the file name and the extension to put it in the template

type Asset struct {
	Name      string
	Extension FileExtension
	Path      string
}

type FileExtension string

const (
	Javascript FileExtension = "js"
	Stylesheet FileExtension = "css"
)

func (dojo *Dojo) Assets() []Asset {
	var assets []Asset

	files, err := os.ReadDir(dojo.Configuration.Assets.Path)
	if err != nil {
		dojo.Logger.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			dojo.Logger.Debugf("assets: skip dir %s", file.Name())
			continue
		}

		dojo.Logger.Debugf("assets: register file %s", file.Name())
		parts := strings.Split(file.Name(), ".")
		if len(parts) > 1 {
			assets = append(assets, Asset{
				Name:      parts[0],
				Extension: FileExtension(parts[1]),
				Path:      fmt.Sprintf("%s/assets/%s", dojo.Configuration.App.Domain, file.Name()),
			})
		}
	}

	return assets
}
