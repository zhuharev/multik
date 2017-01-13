package admin

import (
	"fmt"
	"strings"

	"github.com/Unknwon/com"

	//"pure/multik/models"
	"pure/multik/modules/middleware"
)

func MenuList(c *middleware.Context) {
	// todo handle error
	list, _ := c.E.MenuList()
	c.Data["list"] = list
	c.HTML(200, "admin/menus/list")
}

func MenuCreate(c *middleware.Context) {
	if title := c.Query("title"); title != "" {
		c.E.MenuCreate(title, c.Query("slug"))
	}

	c.Redirect("/admin/menus")
}

func MenuEdit(c *middleware.Context) {
	// todo error
	menu, _ := c.E.MenuGet(c.QueryInt64("id"))
	c.Data["menu_for_edit"] = menu
	c.HTML(200, "admin/menus/edit")
}

func MenuItemCreate(c *middleware.Context) {
	var (
		title  = c.Query("title")
		link   = c.Query("link")
		menuId = c.QueryInt64("menu_id")
		parent = c.QueryInt64("parent")
	)
	if link != "" && title != "" && menuId != 0 {
		c.E.MenuItemCreate(title, link, menuId, parent)
	}

	c.Redirect("/admin/menus/edit?id=" + fmt.Sprint(menuId))
}

func MenuSetPosition(c *middleware.Context) {
	var (
		pos    []int64
		posStr = strings.Split(c.Query("positions"), ";")
	)

	fmt.Println(posStr)
	for _, v := range posStr {
		pos = append(pos, com.StrTo(v).MustInt64())
	}
	e := c.E.MenuSetPosition(pos)
	if e != nil {
		fmt.Println(e)
	}
	c.Redirect("/admin/menus/edit?id=" + c.Query("id"))
}

func MenuItemEdit(c *middleware.Context) {
	var (
		id = c.QueryInt64("id")
	)
	mi, e := c.E.MenuItemGet(id)
	if e != nil {
		fmt.Println(e)
		c.Redirect("/admin/menus")
		return
	}
	if c.Req.Method == "POST" {
		mi.Id = id
		mi.Title = c.Query("title")
		mi.Link = c.Query("link")
		e = c.E.Save(mi)
	}
	if e != nil {
		fmt.Println(e)
		c.Redirect("/admin/menus")
		return
	}
	c.Data["menuitem"] = mi
	c.HTML(200, "admin/menus/itemedit")

}
