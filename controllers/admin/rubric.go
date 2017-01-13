package admin

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/zhuharev/object"
	"pure/multik/models"
	"pure/multik/modules/middleware"
)

func RubricCreate(c *middleware.Context) {
	if c.Req.Method == "POST" {
		rub, e := c.E.RubricsCreate(c.Query("title"),
			c.Query("slug"),
			c.QueryInt64("parent"),
			object.Object(c.QueryInt("item_type")))
		if e != nil {
			c.Flash.Error(e.Error())
			c.Redirect("/admin/rubric/add")
			return
		}

		rubs, e := c.E.RubricsGetAll()
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			c.Redirect("/admin/rubric/add")
			return
		}

		_, e = c.E.RouteSave(rubs.PathString(rub.Id), object.Object(c.QueryInt("item_type")), rub.Id)
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			c.Redirect("/admin/rubric/add")
			return
		}

		c.Flash.Success("Рубрика добавлена")
		c.Redirect("/admin/rubric/" + fmt.Sprint(rub.Id) + "/edit")
		return
	}
	c.HTML(200, "admin/rubric/create")

}

func RubricList(c *middleware.Context) {

	c.HTML(200, "admin/rubric/list")
}

func RubricEdit(c *middleware.Context) {
	rub, e := c.E.RubricsGetById(c.ParamsInt64(":id"))
	if e != nil {
		c.Flash.Error(e.Error())
		c.Redirect("/admin/rubric")
		return
	}
	c.Data["rubric"] = rub
	c.HTML(200, "admin/rubric/edit")
}

func RubricDelete(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	rubric := new(models.Rubric)
	rubric.Id = id
	_, e := c.E.Delete(rubric)
	if e != nil {
		c.Flash.Error(e.Error())
	} else {
		c.Flash.Success("Рубрика удалена")
	}
	c.Redirect("/admin/rubric")
}
