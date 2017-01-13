package models

import (
	"github.com/fatih/color"
	//"github.com/go-xorm/xorm"
	"github.com/zhuharev/eav"
)

func EntityValueId(x *Engine, attribId int64, value string) (id int64, e error) {
	var (
		res eav.Value
	)
	has, e := x.Where("value = ? and attribute_id = ?", value, attribId).Get(&res)
	if !has {
		res.Value = value
		res.AttributeId = attribId
		_, e = x.Insert(&res)
		if e != nil {
			return
		}
	}
	id = res.Id
	return
}

func (x *Engine) EntityValueId(attribId int64, value string) (id int64, e error) {
	return EntityValueId(x, attribId, value)
}

func EntityGetValues(x *Engine, ids ...int64) ([]*eav.Value, error) {
	if ids == nil {
		return nil, nil
	}
	var (
		res = []*eav.Value{}
	)
	e := x.In("id", ids).Find(&res)
	return res, e
}

func (x *Engine) EntityGetValues(ids ...int64) ([]*eav.Value, error) {
	return EntityGetValues(x, ids...)
}

func EntityGetAttributes(x *Engine, entityId int64) (eav.Attributes, error) {

	var (
		attribs = []*eav.Attribute{}
		resA    = []*eav.Value{}
		res     = []int64{}
	)

	ent, e := EntityGet(x, entityId)
	if e != nil {
		return nil, e
	}

	e = x.In("id", ent.ValueIds).Distinct("attribute_id").Find(&resA)
	if e != nil {
		return nil, e
	}

	for _, v := range resA {
		color.Green("%v", v)
		res = append(res, v.AttributeId)
	}

	e = x.In("id", res).Find(&attribs)
	if e != nil {
		return nil, e
	}

	return eav.Attributes(attribs), nil
}

func (x *Engine) EntityGetAttributes(entityId int64) (eav.Attributes, error) {
	return EntityGetAttributes(x, entityId)
}

func EntityAddAttribyte(x *Engine, name string, typ eav.Type) (*eav.Attribute, error) {
	attr := new(eav.Attribute)
	attr.Name = name
	attr.ValueType = typ

	_, e := x.Insert(attr)
	return attr, e
}

func (x *Engine) EntityAddAttribyte(name string, typ eav.Type) (*eav.Attribute, error) {
	return EntityAddAttribyte(x, name, typ)
}

func EntityAllAttribytes(x *Engine) ([]*eav.Attribute, error) {
	var (
		attrs = []*eav.Attribute{}
	)
	e := x.Find(&attrs)
	return attrs, e
}

func (x *Engine) EntityAllAttribytes() ([]*eav.Attribute, error) {
	return EntityAllAttribytes(x)
}
