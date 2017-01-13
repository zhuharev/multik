package controllers

import (
	"os"

	"github.com/fatih/color"
	"github.com/zhuharev/object"
	"pure/multik/models"
	"pure/multik/modules/middleware"
	"pure/multik/modules/setting"
	"pure/multik/modules/widget"

	"pure/multik/controllers/entity"
	"pure/multik/controllers/posts"
	"pure/multik/controllers/rubric"
)

func GlobalInit() {
	setting.NewContext()
	models.NewContext()
	e := widget.NewContext()
	if e != nil {
		color.Red("%s", e.Error())
		os.Exit(1)
	}
}

func Index(c *middleware.Context) {
	c.HTML(200, "index")
}

func Page(c *middleware.Context) {

	p, e := c.E.PostGet(c.QueryInt64("id"))
	if e != nil {
		color.Red("%s", e)
	}

	c.Data["page"] = p

	c.HTML(200, "page")
}

/*func Rubric(c *middleware.Context) {
	color.Green("Navigated to %s", c.Req.Request.URL.Path)

	itemType, e := c.E.GetItemTypeFromUrl(c.Req.Request.URL)
	if e != nil {
		color.Red("%s", e)
	}

	color.Green("%d", itemType)

	switch itemType {
	case object.Entity:
		entity.Item(c)
		return
	case object.List:
		rubric.Rubric(c)
		return
	case object.Post:
		posts.Show(c)
		return
	}

	c.HTML(200, "page")
}*/

func Router(c *middleware.Context) {
	route, e := c.E.RouteGet(c.Req.URL.Path)
	if e != nil {
		c.Flash.Error(e.Error())
		color.Red("%s", e)
		return
	}

	color.Green("Route found %v", route)

	switch route.ObjectType {
	case object.Post:
		posts.Show(c)
		return
	case object.Entity:
		entity.Item(c)
		return
	case object.EntityList, object.PostList, object.List:
		rubric.Rub(c, route.ObjectType, route.ObjectId)
		return
	}
}
