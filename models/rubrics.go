package models

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/zhuharev/object"
)

type Rubric struct {
	Id    int64
	Title string
	Slug  string `xorm:"unique index"`

	ParentId int64

	ChildType object.Object

	Parent *Rubric   `xorm:"-"`
	Childs []*Rubric `xorm:"-"`

	path string `xorm:"-"`
}

func (r Rubric) Path() string {
	return r.path
}

type Rubrics []*Rubric

func (r Rubrics) GetAllParents(parentId int64) (Rubrics, error) {
	var (
		res Rubrics
		lim = 99
	)

	for {
		for _, v := range r {
			if v.Id == parentId {
				res = append(res, v)
				parentId = v.ParentId
			}
		}
		if parentId == 0 || lim < 0 {
			break
		}
		lim--
	}
	return res, nil
}

func (r Rubrics) GetAllChilds(parentId int64) (Rubrics, error) {
	var (
		res Rubrics
	)

	for _, v := range r {
		if v.ParentId == parentId {
			res = append(res, v)
		}
	}
	return res, nil
}

func (r Rubrics) GetAllForType(t string) (res Rubrics) {
	var (
		tt object.Object
	)

	switch t {
	case "entity":
		tt = object.EntityList
	case "post":
		tt = object.PostList
	case "rubric":
		tt = object.List
	}

	for _, v := range r {
		if v.ChildType == tt {
			res = append(res, v)
		}
	}

	return
}

func (r Rubrics) PathString(parentId int64) string {

	var (
		parents, _ = r.GetAllParents(parentId)
		path       = "/"
	)

	for _, v := range parents {
		path = filepath.Join("/", v.Slug+path)
	}

	return path
}

func (r Rubrics) GetById(id int64) *Rubric {
	for _, v := range r {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (r Rubrics) WithPath(rubricId int64) *Rubric {
	rub := r.GetById(rubricId)
	if rub.path != "" {
		return rub
	}
	rub.path = r.PathString(rub.Id)
	return rub
}

func RubricsGetAll(x *Engine) (Rubrics, error) {
	var (
		res Rubrics
	)

	e := x.Find(&res)

	for i, v := range res {
		path := res.PathString(v.Id)
		res[i].path = path
	}

	return res, e
}

func (x *Engine) RubricsGetAll() (Rubrics, error) {
	return RubricsGetAll(x)
}

func RubricsCreate(x *Engine, title, slug string, parent int64, childType object.Object) (*Rubric, error) {
	r := new(Rubric)
	r.Title = title
	r.Slug = slug
	r.ParentId = parent
	r.ChildType = childType

	_, e := x.Insert(r)

	return r, e
}

func (x *Engine) RubricsCreate(title, slug string, parent int64, childType object.Object) (*Rubric, error) {
	return RubricsCreate(x, title, slug, parent, childType)
}

func RubricsGetBySlug(x *Engine, slug string) (*Rubric, error) {
	var (
		rub = new(Rubric)
	)

	_, e := x.Where("slug = ?", slug).Get(rub)

	return rub, e
}

func (x *Engine) RubricsGetBySlug(slug string) (*Rubric, error) {
	return RubricsGetBySlug(x, slug)
}

func RubricsGetById(x *Engine, id int64) (*Rubric, error) {
	var (
		res = new(Rubric)
	)
	_, e := x.Id(id).Get(res)
	return res, e
}

func (x *Engine) RubricsGetById(id int64) (*Rubric, error) {
	return RubricsGetById(x, id)
}

func GetSlugFromUrl(rawurl string) (string, error) {

	u, e := url.Parse(rawurl)
	if e != nil {
		return "", e
	}

	pathArr := strings.Split(u.Path, "/")

	hasId := false
	if id := getIdFromSlug(u.Path); id != 0 {
		hasId = true
	}
	if !hasId {
		return pathArr[len(pathArr)-1], nil
	} else {
		return pathArr[len(pathArr)-2], nil
	}
	return "", nil
}

func GetRubricType(x *Engine, u *url.URL) (object.Object, error) {
	pathArr := strings.Split(u.Path, "/")
	rubSlug := pathArr[len(pathArr)-1]

	println(rubSlug)

	rub, e := x.RubricsGetBySlug(rubSlug)
	if e != nil {
		return 0, e
	}

	return rub.ChildType, nil
}

func (x *Engine) GetRubricType(u *url.URL) (object.Object, error) {
	return GetRubricType(x, u)
}

func GetItemTypeFromUrl(x *Engine, u *url.URL) (object.Object, error) {
	pathArr := strings.Split(u.Path, "/")

	// if not has id, then path is rubric list
	hasId := false
	if id := getIdFromSlug(u.Path); id != 0 {
		hasId = true
	}
	if !hasId {
		post, e := x.PostGetBySlug(u.Path)
		if e != nil {
			return 0, e
		}
		if post.Id != 0 {
			return object.Post, nil
		}
		return object.List, nil
	} else {
		rubSlug := pathArr[len(pathArr)-2]

		rub, e := x.RubricsGetBySlug(rubSlug)
		if e != nil {
			return 0, e
		}
		return rub.ChildType, nil
	}
	return 0, nil
}

func (x *Engine) GetItemTypeFromUrl(u *url.URL) (object.Object, error) {
	return GetItemTypeFromUrl(x, u)
}
