package models

import (
	"github.com/go-xorm/xorm"
)

type Engine struct {
	*xorm.Engine
	site string
}

func (x *Engine) Site() string {
	return x.site
}
