package admin

import (
	"pure/multik/modules/middleware"

	"github.com/fatih/color"
	"github.com/zhuharev/eav"
)

// EavAttributeAdd add attribute controller
func EavAttributeAdd(c *middleware.Context) {
	if c.Req.Method == "POST" {
		attr, e := c.E.EntityAddAttribyte(c.Query("name"), eav.Type(c.QueryInt("type")))
		if e != nil {
			color.Red("%s", e)
			return
		}
		color.Green("%v", attr)
	}

	c.HTML(200, "admin/eav/attribute_add")
}
