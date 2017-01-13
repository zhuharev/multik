package models

import (
	"time"

	"pure/multik/modules/filedb"
)

type File struct {
	Id          int64
	Title       string
	Description string

	Mime string

	FileId int64

	UploadedBy int64
	Created    time.Time `xorm:"created"`
}

func (x *Engine) Upload(f *File, data []byte) error {
	fid, e := filedb.Set(x.Site(), data)
	if e != nil {
		return e
	}

	if f == nil {
		f = new(File)
	}

	f.FileId = fid

	return x.Save(f)
}

func (x *Engine) FilesGet(id int64) (*File, error) {
	f := new(File)
	_, e := x.Id(id).Get(f)
	return f, e
}

func (x *Engine) FilesList() ([]*File, error) {
	var (
		res []*File
	)
	e := x.Find(&res)

	return res, e
}
