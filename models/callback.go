package models

import (
	"time"
)

type Callback struct {
	Id int64

	Name  string
	Phone string
	Text  string

	Processed bool

	Created time.Time `xorm:"created"`
}

func (eng *Engine) CallbackSave(cb *Callback) error {
	if cb.Id == 0 {
		_, e := eng.Insert(cb)
		return e
	} else {
		_, e := eng.Id(cb.Id).Update(cb)
		return e
	}
}

func (eng *Engine) CallbackGet(id int64) (*Callback, error) {
	var (
		cb = new(Callback)
	)
	_, e := eng.Id(id).Get(&cb)
	return cb, e
}

func (eng *Engine) CallbackFind(onlyNotProcessed bool, offset, limit int) ([]*Callback, error) {
	var (
		res []*Callback
	)
	sess := eng.Limit(limit, offset)
	if onlyNotProcessed {
		sess = sess.Where("processed = ?", false)
	}
	e := sess.OrderBy("id desc").Find(&res)
	return res, e
}
