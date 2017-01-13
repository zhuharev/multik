package wiki

import (
	"fmt"

	"pure/multik/models"
	"pure/multik/modules/middleware"
)

func Edit(c *middleware.Context) {
	if c.Req.Method == "POST" {
		slug := c.Params(":slug")
		body := c.Query("body")

		e := models.WikiSave(c.E, slug, []byte(body))
		if e != nil {
			fmt.Println(e)
		}
		c.Redirect("/wiki/" + slug)
		return
	}

	page, e := models.WikiGetBySlug(c.E, c.Params(":slug"))
	if e != nil {
		fmt.Println(e)

	}

	c.Data["WikiPage"] = page
	c.HTML(200, "wiki/edit")
}
