package models

import (
	"github.com/zhuharev/object"
)

type rubricItemView struct {
	slug  string
	title string
	typ   object.Object
	image string
}

func (riv rubricItemView) Slug() string {
	return riv.slug
}

func (riv rubricItemView) Title() string {
	return riv.title
}
func (riv rubricItemView) Type() object.Object {
	return riv.typ
}
func (riv rubricItemView) Image() string {
	return riv.image
}

// convert Entity to rubric view
func (ent Entity) ToRubricView() RubricsItemView {
	return rubricItemView{title: ent.Title, typ: object.Entity, slug: ent.Slug()}
}

func (rub Rubric) ToRubricView() RubricsItemView {
	return rubricItemView{title: rub.Title, typ: object.List, slug: rub.Slug}
}
