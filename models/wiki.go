package models

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

type WikiPage struct {
	Id   int64
	Slug string `xorm:"unique index"`

	Description string
	Body        []byte

	PrevDelta int64
	Owner     int64
	Created   int64
	Deleted   int64
	Updated   int64
	Version   int64
}

func (w WikiPage) BodyString() string {
	return string(w.Body)
}

type WikiDelta struct {
	Id          int64
	ItemId      int64
	PrevDelta   int64
	Description string

	Delta   []byte
	Owner   int64
	Created int64
}

func WikiGetBySlug(x *Engine, slug string) (*WikiPage, error) {
	w := &WikiPage{}
	has, e := x.Where("slug = ?", slug).Get(w)
	if !has {
		return nil, ErrNotFound
	}
	return w, e
}

func WikiSave(x *Engine, slug string, body []byte) error {
	wOld, e := WikiGetBySlug(x, slug)
	if e != nil {
		if e != ErrNotFound {
			return e
		} else {
			return WikiCreate(x, slug, body)
		}
	}
	textOld := string(wOld.Body)
	textNew := string(body)

	d := diffmatchpatch.New()
	b := d.DiffMain(textNew, textOld, false)

	dl := d.DiffToDelta(b)

	delta := &WikiDelta{
		ItemId:    wOld.Id,
		PrevDelta: wOld.PrevDelta,
		Delta:     []byte(dl),
	}

	e = x.Save(delta)
	if e != nil {
		return e
	}

	w := &WikiPage{Body: body, PrevDelta: delta.Id}
	_, e = x.Cols("prev_delta", "body").Id(wOld.Id).Update(w)

	if e != nil {
		return e
	}
	return nil
}

func WikiCreate(x *Engine, slug string, body []byte) error {
	w := new(WikiPage)
	w.Slug = slug
	w.Body = body

	_, e := x.Insert(w)
	return e
}

func WikiHistoryBySlug(x *Engine, slug string) ([]*WikiDelta, error) {
	var wd []*WikiDelta
	e := x.Where("item_id = (select id from wiki_page where slug = ?)", slug).OrderBy("id desc").Find(&wd)
	return wd, e
}

func WikiHistory(x *Engine, id int64) ([]*WikiDelta, error) {
	var wd []*WikiDelta
	e := x.Where("item_id = ?", id).OrderBy("id desc").Find(&wd)
	return wd, e
}
