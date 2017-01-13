package admin

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/zhuharev/object"
	//"pure/multik/models"
	"pure/multik/modules/middleware"
)

func EntityList(c *middleware.Context) {
	var (
		p           = c.QueryInt("p")
		itemsInPage = 10
	)
	if p < 1 {
		p = 1
	}
	items, e := c.E.EntityList(p, itemsInPage)
	if e != nil {
		c.Error(500, e.Error())
		return
	}
	c.Data["items"] = items
	c.HTML(200, "admin/entity/list")
}

func EntityAdd(c *middleware.Context) {
	if c.Req.Method == "POST" {
		ent, e := c.E.EntityAdd(c.Query("title"))
		if e != nil {
			color.Red("%s", e)
		}
		color.Green("added %v", ent)
		c.Flash.Success("Вы успешно вошли")
		c.Redirect("/admin/entity/" + fmt.Sprint(ent.Id) + "/edit")
		return
	}

	c.HTML(200, "admin/entity/add")
}

func EntityEdit(c *middleware.Context) {
	ent, e := c.E.EntityGet(c.ParamsInt64(":id"))
	if e != nil {
		color.Red("%s", e)
	}

	entAttrs, e := c.E.EntityGetAttributes(ent.Id)
	if e != nil {
		color.Red("%s", e)
	}

	color.Green("%v", ent)
	if c.Req.Method == "POST" {
		ent.Title = c.Query("title")
		ent.SetPrice(c.QueryFloat64("price"))
		color.Yellow("%v %v", ent.Price(), c.QueryFloat64("price"))
		e := c.E.EntityUpdate(ent)
		if e != nil {
			color.Red("%s", e)
		}

		/* rubric */
		rubricId := c.QueryInt64("rubric")
		_, e = c.E.RubricItemSave(rubricId, object.Entity, ent.Id)
		if e != nil {
			color.Red("%s", e)
		}

		rubs, e := c.E.RubricsGetAll()
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			return
		}

		_, e = c.E.RouteSave(filepath.Join(rubs.PathString(rubricId), ent.Slug()), object.Entity, ent.Id)
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			c.Redirect("/admin/dasboard")
			return
		}

		ent.RubricId = rubricId
	}

	attrs, e := c.E.EntityAllAttribytes()
	if e != nil {
		color.Red("%s", e)
	}
	c.Data["entAttrs"] = entAttrs
	c.Data["attributes"] = attrs
	c.Data["entity"] = ent
	c.HTML(200, "admin/entity/edit")
}

func EntityValueAdd(c *middleware.Context) {
	ent, e := c.E.EntityGet(c.QueryInt64("id"))
	if e != nil {
		color.Red("%s", e)
	}

	id, e := c.E.EntityValueId(c.QueryInt64("attribute_id"), c.Query("value"))
	if e != nil {
		color.Red("%s", e)
	}

	ent.ValueIds = append(ent.ValueIds, id)
	e = c.E.EntityUpdate(ent)
	if e != nil {
		color.Red("%s", e)
	}
	c.Redirect("/admin/entity/" + c.Query("id") + "/edit")
}

func EntityUpload(c *middleware.Context) {
	var (
		entId = c.ParamsInt64(":id")
	)
	imgId, e := upload(c)
	if e != nil {
		color.Red("%s", e.Error())
	}

	e = c.E.EntityAddImage(entId, imgId)
	if e != nil {
		color.Red("%s", e.Error())
	}
}
