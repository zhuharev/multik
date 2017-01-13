package admin

import (
	"pure/multik/modules/middleware"
	//"github.com/fatih/color"
)

func Dashboard(c *middleware.Context) {

	c.HTML(200, "admin/dashboard")
}
