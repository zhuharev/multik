package models

import (
	"github.com/Unknwon/com"
	//"fmt"
	"github.com/zhuharev/object"
	"time"
	//"github.com/sisteamnik/guseful/chpu"
)

type Post struct {
	Id    int64
	Title string
	Slug  string `xorm:"unique index"`
	Body  string

	RubricId int64

	Created time.Time `xorm:"created"`
}

type postRubricView struct {
	id    int64
	title string
	slug  string
	typ   object.Object
	image string
}

func (ent postRubricView) Slug() string {
	return ent.slug
}

func (ent postRubricView) Title() string {
	return ent.title
}
func (ent postRubricView) Type() object.Object {
	return ent.typ
}
func (ent postRubricView) Image() string {
	return ent.image
}

func (p Post) ToRubricView() RubricsItemView {
	return postRubricView{
		title: p.Title,
		image: "",
		typ:   object.Post,
		id:    p.Id,
		slug:  p.Slug,
	}
}

func (x *Engine) PostSave(id, rubricId int64, title, slug, body string) (*Post, error) {
	p := new(Post)
	p.Id = id
	p.RubricId = rubricId
	p.Title = title
	p.Body = body
	p.Slug = slug
	e := x.Save(p)
	return p, e
}

func (e *Engine) PostGet(id int64) (*Post, error) {
	return getPost(e, id)
}

func getPost(x *Engine, id int64) (*Post, error) {
	var (
		p = new(Post)
		e error
	)
	_, e = x.Id(id).Get(p)
	return p, e
}

func (x *Engine) PostGetBySlug(slug string) (*Post, error) {
	var (
		p = new(Post)
		e error
	)
	_, e = x.Where("slug = ?", slug).Get(p)
	return p, e
}

func (x *Engine) PostList(limit int, offsets ...int) ([]*Post, error) {
	var (
		res []*Post
	)
	e := x.Limit(limit, offsets...).OrderBy("id asc").Find(&res)
	return res, e
}

func (x *Engine) PostRubricList(rubricId int64, limitIf interface{}, offsets ...int) ([]*Post, error) {
	var (
		res   []*Post
		limit int
	)
	switch limitIf.(type) {
	case int64:
		limit = int(limitIf.(int64))
	case string:
		limit = com.StrTo(limitIf.(string)).MustInt()
	case int:
		limit = limitIf.(int)
	}
	e := x.Where("rubric_id = ?", rubricId).Limit(limit, offsets...).OrderBy("id asc").Find(&res)
	return res, e
}
