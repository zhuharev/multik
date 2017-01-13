package admin

import (
	"pure/multik/modules/middleware"

	"github.com/fatih/color"
)

func Login(c *middleware.Context) {
	if c.Req.Method == "POST" {
		if c.Query("name") == "admin" && c.Query("pass") == "admin" {
			e := c.Sess.Set("is_admin", true)
			if e != nil {
				color.Red("%s", e)
				c.Redirect("/")
				return
			}
			c.Flash.Success("Вы успешно вошли")
			c.Redirect("/admin/dashboard")
			return
		}
		c.Flash.Error("Логин или пароль не верны")
		c.Redirect("/admin/login")
		return
	}

	c.HTML(200, "admin/login")
}
