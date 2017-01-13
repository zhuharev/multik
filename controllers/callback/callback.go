package callback

import (
	"pure/multik/models"
	"pure/multik/modules/middleware"

	"log"
)

func New(c *middleware.Context) {
	if c.Req.Method == "POST" {
		cb := &models.Callback{
			Name:  c.Query("name"),
			Phone: c.Query("phone"),
			Text:  c.Query("text"),
		}
		e := c.E.CallbackSave(cb)
		if e != nil {
			log.Println(e)
		}
	}
	c.HTML(200, "callback")
}
