package admin

import (
	//"fmt"

	"github.com/fatih/color"
	//"pure/multik/models"
	"pure/multik/modules/middleware"
)

func WidgetEdit(c *middleware.Context) {

	key := "widget." + c.Params(":slug") + "." + "data"

	if c.Req.Method == "POST" {
		value := c.Query("widget_setting")

		e := c.E.Set(key, []byte(value))
		if e != nil {
			color.Red("%s", e)
			c.Flash.Error(e.Error())
			c.Redirect(c.Req.Request.URL.Path)
			return
		}
		c.Flash.Success("Сохранено")
		c.Redirect(c.Req.Request.URL.Path)
		return
	}

	data, e := c.E.Get(key)
	if e != nil {
		color.Red("%s", e)
		c.Flash.Error(e.Error())
		c.Redirect(c.Req.Request.URL.Path)
		return
	}
	c.Data["widget_setting"] = string(data)

	c.HTML(200, "admin/widget/edit")
}
