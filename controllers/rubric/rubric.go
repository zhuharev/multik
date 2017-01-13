package rubric

import (
	"github.com/fatih/color"
	"github.com/zhuharev/object"

	//"pure/multik/models"
	"pure/multik/modules/middleware"
)

func Rub(c *middleware.Context, ot object.Object, rubricId int64) {
	color.Green("Get item for %s", ot)

	items, e := c.E.RubricsGetItemsById(rubricId, ot)
	if e != nil {
		panic(e)
	}

	rub, e := c.E.RubricsGetById(rubricId)
	if e != nil {
		panic(e)
	}

	c.Data["currentRubric"] = rub
	c.Data["items"] = items
	c.HTML(200, "rubric/list")
}

/*func Rubric(c *middleware.Context) {

	slug, e := models.GetSlugFromUrl(c.Req.Request.URL.String())
	if e != nil {
		panic(e)
	}

	typ, e := c.E.GetRubricType(c.Req.Request.URL)
	if e != nil {
		panic(e)
	}

	color.Green("[type] %d", typ)

	items, e := c.E.RubricsGetItems(slug, typ)
	if e != nil {
		panic(e)
	}

	rub, e := c.E.RubricsGetBySlug(slug)
	if e != nil {
		panic(e)
	}

	c.Data["currentRubric"] = rub
	c.Data["items"] = items

	list := "list"

	switch typ {
	//	case models.RI_POST:
	//		list = "list_post"
	//	case models.RI_ENTITY:
	//		list = "list_entity"
	//	case models.RI_RUBRIC:
	//		list = "list_rubric"
	}

	c.HTML(200, "rubric/"+list)
}*/
