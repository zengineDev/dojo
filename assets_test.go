package dojo

import "testing"

func TestApplication_Assets(t *testing.T) {
	app := New(DefaultConfiguration{
		Assets: AssetsConfigs{Path: "./_test/dist"},
	})

	assets := app.Assets()
	if len(assets) == 0 {
		t.Error("no files")
	}

	for _, a := range assets {
		if a.Name != "index" {
			t.Errorf("filename: %s dose not match", a.Name)
		}
	}
}
