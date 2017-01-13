package widget

import (
	//"html/template"
	"encoding/json"
	"os"
	"path/filepath"

	//"github.com/fatih/color"
	"github.com/fatih/color"
	"github.com/ungerik/go-dry"
	//"github.com/zhuharev/boltutils"
	"pure/multik/modules/base"
)

var (
	widgets = map[string][]*Widget{}
)

type Widget struct {
	Name string       `json:"name"`
	Slug string       `json:"slug"`
	Data []WidgetData `json:"data"`
}

type WidgetData struct {
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func Widgets(site string) []*Widget {
	if widgets != nil {
		return widgets[site]
	}
	return nil
}

func NewContext() error {
	sites, e := base.Sites()
	if e != nil {
		return e
	}

	for _, site := range sites {

		siteWidgets := []*Widget{}

		dirPath := "sites/" + site + "/widgets/"
		if !dry.FileExists(dirPath) {
			e = os.MkdirAll(dirPath, 0777)
			if e != nil {
				return e
			}
		}

		list, e := dry.ListDirDirectories(dirPath)
		if e != nil {
			return e
		}

		for _, widget := range list {
			w := new(Widget)
			bts, e := dry.FileGetBytes(filepath.Join(dirPath, widget, "manifest.json"))
			if e != nil {
				return e
			}
			e = json.Unmarshal(bts, w)
			if e != nil {
				return e
			}
			color.Green("Found widget %s for site %s", widget, site)
			siteWidgets = append(siteWidgets, w)
		}

		widgets[site] = siteWidgets

	}
	return nil
}
