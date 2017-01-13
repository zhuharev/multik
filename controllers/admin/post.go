package admin

import (
	"pure/multik/modules/middleware"

	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/zhuharev/object"
)

func PostEdit(c *middleware.Context) {
	var (
		id = c.QueryInt64("id")
	)
	if c.Req.Method == "POST" {
		p, e := c.E.PostSave(id, c.QueryInt64("rubric"),
			c.Query("title"), c.Query("slug"), c.Query("body"))
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			c.Redirect("/admin/posts/edit?id=" + fmt.Sprint(p.Id))
			return
		}

		rubs, e := c.E.RubricsGetAll()
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			c.Redirect("/admin/posts/edit?id=" + fmt.Sprint(p.Id))
			return
		}

		_, e = c.E.RouteSave(filepath.Join(rubs.PathString(p.RubricId), p.Slug), object.Post, p.Id)
		if e != nil {
			color.Red("%s", e.Error())
			c.Flash.Error(e.Error())
			c.Redirect("/admin/posts/edit?id=" + fmt.Sprint(p.Id))
			return
		}
		c.Redirect("/admin/posts/edit?id=" + fmt.Sprint(p.Id))
		return
	}

	post, e := c.E.PostGet(id)
	if e != nil {
		panic(e)
	}

	c.Data["post"] = post
	c.HTML(200, "admin/post/edit")
}
