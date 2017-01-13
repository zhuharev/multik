package admin

import (
	//"pure/multik/models"
	"pure/multik/modules/middleware"

	"log"
)

func CallbackFind(c *middleware.Context) {
	callbacks, e := c.E.CallbackFind(false, 0, 20)
	if e != nil {
		log.Println(e)
	}
	c.Data["callbacks"] = callbacks
	c.HTML(200, "admin/callback/find")
}
