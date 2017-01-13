package models

import (
	"sort"
	//"strings"
	"time"

	"path/filepath"
)

type Menu struct {
	Id    int64
	Title string
	Slug  string

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`

	Items MenuItems `xorm:"-"`
}

func MenuGetAll(x *Engine) ([]*Menu, error) {
	var menus []*Menu
	e := x.Find(&menus)
	if e != nil {
		return nil, e
	}
	return menus, e
}

func (x *Engine) MenuGetAll() ([]*Menu, error) {
	return MenuGetAll(x)
}

type MenuItem struct {
	Id    int64
	Title string
	Link  string

	MenuId   int64
	ParentId int64

	Position int

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type MenuItems []MenuItem

func (a MenuItems) Len() int           { return len(a) }
func (a MenuItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MenuItems) Less(i, j int) bool { return a[i].Position < a[j].Position }

func (m MenuItems) FirstLevel() (res MenuItems) {
	for _, v := range m {
		if v.ParentId == 0 {
			res = append(res, v)
		}
	}
	sort.Sort(res)
	return
}

func (m MenuItems) HasChilds(id int64) bool {
	for _, v := range m {
		if v.ParentId == id {
			return true
		}
	}
	return false
}

func (m MenuItems) Childs(id int64) (res MenuItems) {
	for _, v := range m {
		if v.ParentId == id {
			//if strings.HasPrefix(s, prefix)
			res = append(res, v)
		}
	}
	sort.Sort(res)
	return
}

func (r MenuItems) GetAllParents(parentId int64) (MenuItems, error) {
	var (
		res MenuItems
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

func (m MenuItems) Path(id int64) string {

	var (
		item MenuItem
	)

	for i := range m {
		if m[i].Id == id {
			item = m[i]
		}
	}

	var (
		parents, _ = m.GetAllParents(item.ParentId)
		path       = "/"
	)

	for _, v := range parents {
		path = filepath.Join(v.Link, path)
	}

	path = filepath.Join(path, item.Link)

	return path
}

func (m MenuItems) Menu(id int64) (items MenuItems) {
	for _, v := range m {
		if v.MenuId == id {
			items = append(items, v)
		}
	}
	return
}

func MenuCreate(x *Engine, title string, slug string) (*Menu, error) {
	m := &Menu{
		Title: title,
		Slug:  slug,
	}
	e := x.Save(m)
	return m, e
}

func (x *Engine) MenuCreate(title string, slug string) (*Menu, error) {
	return MenuCreate(x, title, slug)
}

func MenuItemCreate(x *Engine, title, link string, menuId int64, parent int64) (*MenuItem, error) {
	m := &MenuItem{
		Title:    title,
		Link:     link,
		MenuId:   menuId,
		ParentId: parent,
	}
	e := x.Save(m)
	return m, e
}

func (x *Engine) MenuItemCreate(title, link string, menuId int64, parent int64) (*MenuItem, error) {
	return MenuItemCreate(x, title, link, menuId, parent)
}

func MenuItemGet(x *Engine, id int64) (*MenuItem, error) {
	mi := new(MenuItem)
	_, e := x.Id(id).Get(mi)
	return mi, e
}

func (x *Engine) MenuItemGet(id int64) (*MenuItem, error) {
	return MenuItemGet(x, id)
}

func MenuSetPosition(x *Engine, positions []int64) error {
	for position, itemId := range positions {
		item := MenuItem{Position: position + 1}
		_, e := x.Id(itemId).Update(&item)
		if e != nil {
			return e
		}
	}
	return nil
}

func (x *Engine) MenuSetPosition(positions []int64) error {
	return MenuSetPosition(x, positions)
}

func MenuList(x *Engine) ([]*Menu, error) {
	var (
		res []*Menu
	)
	e := x.Find(&res)
	return res, e
}

func (x *Engine) MenuList() ([]*Menu, error) {
	return MenuList(x)
}

func MenuItemsGet(x *Engine, menuId int64) (MenuItems, error) {
	var res []MenuItem
	e := x.Where("menu_id = ?", menuId).Find(&res)
	return res, e
}

func (x *Engine) MenuItemsGet(menuId int64) (MenuItems, error) {
	return MenuItemsGet(x, menuId)
}

func MenuGet(x *Engine, id int64) (*Menu, error) {
	var (
		res = new(Menu)
	)
	_, e := x.Id(id).Get(res)
	if e != nil {
		return res, e
	}
	res.Items, e = MenuItemsGet(x, id)
	return res, e
}

func (x *Engine) MenuGet(id int64) (*Menu, error) {
	return MenuGet(x, id)
}

type Menus struct {
	Menus []*Menu
	Items MenuItems
}

func (m *Menus) Menu(slug string) *Menu {
	for _, v := range m.Menus {
		if v.Slug == slug {
			if v.Items != nil {
				return v
			} else {
				v.Items = m.Items.Menu(v.Id)
				return v
			}
		}
	}
	return &Menu{Items: MenuItems{}}
}

func MenusGet(x *Engine) (*Menus, error) {
	var (
		items = MenuItems{}
		menus = new(Menus)
	)
	e := x.Find(&items)
	if e != nil {
		return menus, e
	}
	menus.Items = items

	allmenu, e := MenuGetAll(x)
	if e != nil {
		return menus, e
	}
	menus.Menus = allmenu

	return menus, nil
}

func (x *Engine) MenusGet() (*Menus, error) {
	return MenusGet(x)
}
