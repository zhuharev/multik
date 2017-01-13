package entity

import (
	"strings"

	"github.com/fatih/color"
	"pure/multik/models"
	"pure/multik/modules/middleware"
)

func Item(c *middleware.Context) {
	arr := strings.Split(c.Req.Request.URL.Path, "/")
	if len(arr) < 2 {
		return

	}
	item(c, arr[len(arr)-1])
}

func item(c *middleware.Context, slug string) {
	var (
		entity *models.Entity
		e      error
	)

	entity, e = c.E.EntityGetBySlug(slug)
	if e != nil {
		//TODO
	}

	color.Green("%v", entity)
	c.Data["result"] = entity
	c.HTML(200, "entity/item")
}
