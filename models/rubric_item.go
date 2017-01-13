package models

import (
	"github.com/fatih/color"
	"github.com/zhuharev/object"
)

type RubricItem struct {
	Id       int64
	RubricId int64
	ItemType object.Object
	ItemId   int64
}

func RubricItemSave(e *Engine, rubricId int64, itemType object.Object, itemId int64) (*RubricItem, error) {
	ri := new(RubricItem)
	has, err := e.Where("item_type = ? and item_id = ?", itemType, itemId).Get(ri)
	if !has {
		if err != nil {
			return nil, err
		}
		ri.ItemId = itemId
		ri.ItemType = itemType
		ri.RubricId = rubricId
		_, err := e.Insert(ri)
		if err != nil {
			return nil, err
		}
	} else {
		if err != nil {
			return nil, err
		}
		if ri.Id != 0 {
			ri.ItemId = itemId
			ri.ItemType = itemType
			ri.RubricId = rubricId
			_, err := e.Id(ri.Id).Update(ri)
			if err != nil {
				return nil, err
			}
		}

	}
	return ri, nil
}

func (e *Engine) RubricItemSave(rubricId int64, itemType object.Object, itemId int64) (*RubricItem, error) {
	return RubricItemSave(e, rubricId, itemType, itemId)
}

func RubricIdGetForItem(e *Engine, itemType object.Object, itemId int64) (int64, error) {
	var (
		ri = new(RubricItem)
	)
	_, err := e.Where("item_type = ? and item_id = ?", itemType, itemId).Get(ri)
	return ri.RubricId, err
}

func (e *Engine) RubricIdGetForItem(itemType object.Object, itemId int64) (int64, error) {
	return RubricIdGetForItem(e, itemType, itemId)
}

type RubricsItemView interface {
	Slug() string
	Title() string
	Image() string
	Type() object.Object
}

func RubricsGetItemsById(x *Engine, id int64, itemType object.Object) ([]RubricsItemView, error) {
	rub, e := x.RubricsGetById(id)
	if e != nil {
		return nil, e
	}

	color.Green("[slug] %s, [child] %d, [itemType] %d, [rubId] %d", rub.Slug, rub.ChildType, itemType, rub.Id)

	var (
		res []RubricsItemView
	)

	switch itemType {
	case object.EntityList:
		var restmp []Entity
		e = x.Where("rubric_id = ?", rub.Id).Find(&restmp)
		if e != nil {
			return nil, e
		}
		res = make([]RubricsItemView, len(restmp))
		for i, v := range restmp {
			res[i] = v.ToRubricView()
		}
		color.Green("getted %d items ", len(res))
		return res, nil
	case object.PostList:
		restmp, e := x.PostRubricList(rub.Id, 10)
		if e != nil {
			return nil, e
		}
		res = make([]RubricsItemView, len(restmp))
		for i, v := range restmp {
			res[i] = v.ToRubricView()
		}
		return res, nil
	case object.List:
		rubs, e := x.RubricsGetAll()
		if e != nil {
			return nil, e
		}
		rubs, e = rubs.GetAllChilds(id)
		if e != nil {
			return nil, e
		}
		res = make([]RubricsItemView, len(rubs))
		for i, v := range rubs {
			res[i] = v.ToRubricView()
		}
		return res, nil
	}

	return nil, nil
}

func (x *Engine) RubricsGetItemsById(id int64, itemType object.Object) ([]RubricsItemView, error) {
	return RubricsGetItemsById(x, id, itemType)
}

/*func RubricsGetItems(x *Engine, slug string, itemType object.Object) ([]RubricsItemView, error) {
	rub, e := x.RubricsGetBySlug(slug)
	if e != nil {
		return nil, e
	}

	color.Green("[slug] %s, [child] %d, [itemType] %d, [rubId] %d", slug, rub.ChildType, itemType, rub.Id)

	var (
		res []RubricsItemView
	)

	switch rub.ChildType {
	case object.EntityList:
		var restmp []Entity
		e = x.Where("rubric_id = ?", rub.Id).Find(&restmp)
		if e != nil {
			return nil, e
		}
		res = make([]RubricsItemView, len(restmp))
		for i, v := range restmp {
			res[i] = v.ToRubricView()
		}
		color.Green("getted %d items ", len(res))
		return res, nil
	case object.PostList:
		restmp, e := x.PostRubricList(rub.Id, 10)
		if e != nil {
			return nil, e
		}
		res = make([]RubricsItemView, len(restmp))
		for i, v := range restmp {
			res[i] = v.ToRubricView()
		}
		return res, nil
	case object.List:
	}

	return nil, nil
}

func (x *Engine) RubricsGetItems(slug string, itemType object.Object) ([]RubricsItemView, error) {
	return RubricsGetItems(x, slug, itemType)
}
*/
