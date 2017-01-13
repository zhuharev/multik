package wiki

import (
	"fmt"

	"pure/multik/models"
	"pure/multik/modules/middleware"
)

func View(c *middleware.Context) {

	page, e := models.WikiGetBySlug(c.E, c.Params(":slug"))
	if e != nil {
		fmt.Println(e)

	}
	c.Data["Title"] = c.Params(":slug")
	c.Data["WikiPage"] = page

	c.HTML(200, "wiki/view")
}
