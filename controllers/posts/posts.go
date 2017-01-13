package posts

import (
	"github.com/fatih/color"
	"strings"

	"pure/multik/modules/middleware"
)

func Show(c *middleware.Context) {

	arr := strings.Split(c.Req.Request.URL.Path, "/")

	color.Yellow("Get post by slug(%s)", arr[len(arr)-1])
	p, e := c.E.PostGetBySlug(arr[len(arr)-1])
	if e != nil {
		color.Red("%s", e)
	}

	c.Data["page"] = p

	c.HTML(200, "page")
}
