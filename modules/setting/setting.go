package setting

import (
	"gopkg.in/ini.v1"
	"path/filepath"
)

var (
	iniFile *ini.File

	Router = struct {
		EntityItem string
	}{}
)

func NewContext() {
	var (
		e error
	)
	iniFile, e = ini.Load(filepath.Join("conf", "app.ini"))
	if e != nil {
		panic(e)
	}

	sec := iniFile.Section("routes")
	Router.EntityItem = sec.Key("entity.item").String()
}
