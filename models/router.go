package models

import (
	"fmt"
	"path/filepath"

	"github.com/zhuharev/object"
)

var (
	ErrRouteNotFound = fmt.Errorf("route not found")
)

type Route struct {
	Id         int64
	Slug       string `xorm:unique index`
	ObjectType object.Object
	ObjectId   int64
}

func (x *Engine) RouteNew(slug string, ot object.Object, oId int64) (*Route, error) {
	var (
		route = new(Route)
	)
	route.Slug = filepath.Clean(slug)
	route.ObjectId = oId
	route.ObjectType = ot
	e := x.Save(route)
	return route, e
}

func (x *Engine) RouteSave(slug string, ot object.Object, oId int64) (*Route, error) {
	var (
		route = new(Route)
	)
	cnt, e := x.Where("object_id = ?", oId).And("object_type = ?", ot).Count(route)
	if e != nil {
		return route, e
	}
	if cnt == 0 {
		return x.RouteNew(slug, ot, oId)
	}

	route.Slug = filepath.Clean(slug)

	_, e = x.Cols("slug").Where("object_id = ?", oId).And("object_type = ?", ot).Update(route)
	return route, e
}

func (x *Engine) RouteGet(slug string) (*Route, error) {
	var (
		route = new(Route)
	)
	slug = filepath.Clean(slug)
	has, e := x.Where("slug = ?", slug).Get(route)
	if e != nil {
		return nil, e
	}
	if !has {
		return nil, ErrRouteNotFound
	}
	return route, nil
}
