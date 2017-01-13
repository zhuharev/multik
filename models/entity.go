package models

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/go-xorm/xorm"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/zhuharev/eav"
	"github.com/zhuharev/object"
)

type Entity struct {
	Id        int64
	Title     string
	PriceData int

	ValueIds   []int64
	Values     []*eav.Value     `db:"-"`
	Attributes []*eav.Attribute `db:"-"`
	RubricId   int64            `db:"-"`

	ManufacturerId int64
	ModelId        int64

	Images []int64

	Uri string
}

func (ent Entity) JSON() ([]byte, error) {
	return json.Marshal(ent)
}

func (ent Entity) MustJSON() string {
	b, _ := json.Marshal(ent)
	return string(b)
}

func (ent Entity) Price() float64 {
	return float64(ent.PriceData) / 100.0
}

func (ent *Entity) SetPrice(price float64) {
	p := int(price * 100)
	color.Yellow("set %v", p)
	ent.PriceData = p
}

func (ent Entity) Slug() string {
	return fmt.Sprintf("%s_%d", chpu.Chpu(ent.Title), ent.Id)
}

func EntityAdd(x *xorm.Engine, title string) (*Entity, error) {
	ent := new(Entity)
	ent.Title = title

	_, e := x.Insert(ent)

	return ent, e
}

func (x *Engine) EntityAdd(title string) (*Entity, error) {
	return EntityAdd(x.Engine, title)
}

func EntityGet(x *Engine, id int64) (*Entity, error) {
	var (
		en = new(Entity)
	)
	_, e := x.Id(id).Get(en)
	if e != nil {
		return en, e
	}

	if en == nil {
		return nil, fmt.Errorf("entity not found")
	}

	values, e := EntityGetValues(x, en.ValueIds...)
	if e != nil {
		return en, e
	}
	en.Values = values

	// rubric
	rubId, e := RubricIdGetForItem(x, object.Entity, id)
	if e != nil {
		return en, e
	}
	en.RubricId = rubId

	return en, e
}

func (x *Engine) EntityGetByIds(ids interface{}) ([]*Entity, error) {

	var (
		res    []*Entity
		idsInt []int64
	)
	switch ids.(type) {
	case string:
		strIds := ids.(string)
		idsStrArr := strings.Split(strIds, ",")
		for _, v := range idsStrArr {
			idsInt = append(idsInt, com.StrTo(v).MustInt64())
		}
	}

	e := x.In("id", idsInt).Find(&res)

	return res, e
}

func (x *Engine) EntityList(page, itemsInPage int) ([]*Entity, error) {
	res := []*Entity{}
	e := x.Limit(itemsInPage, (page-1)*itemsInPage).Find(&res)
	return res, e
}

func (x *Engine) EntityGet(id int64) (*Entity, error) {
	return EntityGet(x, id)
}

func (x *Engine) EntityGetBySlug(slug string) (*Entity, error) {
	return x.EntityGet(getIdFromSlug(slug))
}

func (x *Engine) EntityFullPath(entId int64) (string, error) {

	ent, e := x.EntityGet(entId)
	if e != nil {
		return "", e
	}

	rubId, e := x.RubricIdGetForItem(object.Entity, entId)
	if e != nil {
		return "", e
	}

	rubs, e := x.RubricsGetAll()
	if e != nil {
		return "", e
	}

	rub, e := x.RubricsGetById(rubId)
	if e != nil {
		return "", e
	}
	return path.Join(rubs.PathString(rub.ParentId), rub.Slug, ent.Slug()), nil
}

func EntitySetImages(x *Engine, entId int64, ids []int64) error {
	var (
		ent = new(Entity)
	)
	ent.Images = ids
	_, e := x.Id(entId).Cols("images").Update(ent)
	return e
}

func EntityAddImage(x *Engine, entId int64, imageId int64) error {
	var (
		ent = new(Entity)
	)
	_, e := x.Id(entId).Cols("images").Get(ent)
	if e != nil {
		return e
	}
	ent.Images = append(ent.Images, imageId)
	_, e = x.Id(entId).Cols("images").Update(ent)
	if e != nil {
		return e
	}
	return nil
}
func (x *Engine) EntityAddImage(entId int64, imageId int64) error {
	return EntityAddImage(x, entId, imageId)
}

func EntityUpdate(x *Engine, ent *Entity) error {

	p, e := x.EntityFullPath(ent.Id)
	if e != nil {
		return e
	}

	ent.Uri = p

	_, e = x.RubricItemSave(ent.RubricId, object.Entity, ent.Id)
	if e != nil {
		return e
	}

	_, e = x.Id(ent.Id).Update(ent)
	if e != nil {
		return e
	}

	return e
}

func (x *Engine) EntityUpdate(ent *Entity) error {
	return EntityUpdate(x, ent)
}
