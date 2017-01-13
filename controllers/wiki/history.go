package wiki

import (
	"fmt"

	"pure/multik/models"
	"pure/multik/modules/middleware"
)

func History(c *middleware.Context) {
	deltas, e := models.WikiHistoryBySlug(c.E, c.Params(":slug"))
	if e != nil {
		fmt.Println(e)
	}
	c.Data["Deltas"] = deltas
	c.HTML(200, "wiki/history")
}
