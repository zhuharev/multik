package middleware

import (
	"github.com/fatih/color"
	"github.com/go-macaron/session"
	"github.com/zhuharev/object"
	"gopkg.in/macaron.v1"
	//"html/template"

	"pure/multik/models"
	"pure/multik/modules/widget"
)

type Context struct {
	*macaron.Context

	E *models.Engine

	Sess  session.Store
	Flash *session.Flash
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, flash *session.Flash) {
		ctx := &Context{
			Context: c,
			E:       models.GetEngine(c.Req.Host),
			Sess:    sess,
			Flash:   flash,
		}

		rubs, e := ctx.E.RubricsGetAll()
		if e != nil {
			color.Red("%s", e.Error())
		} else {
			for _, v := range rubs {
				color.Green("%v", v)
			}

			ctx.Data["Rubrics"] = rubs
		}

		menus, e := ctx.E.MenusGet()
		if e != nil {
			color.Red("%s", e.Error())
		} else {
			ctx.Data["Menus"] = menus
		}

		ctx.Data["CurrentUrl"] = ctx.Req.Request.URL.Path

		ctx.Data["AvailableObjects"] = object.Slice

		ctx.Data["widget"] = makeWidget(ctx)
		ctx.Data["Widgets"] = widget.Widgets(ctx.Req.Host)
		c.Map(ctx)
	}
}
