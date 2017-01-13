package filedb

import (
	//"io"
	"path/filepath"

	"github.com/zhuharev/bloblog"
)

func Set(site string, data []byte) (int64, error) {
	bl, e := bloblog.Open(filepath.Join("sites", site, "db", "files.bloblog"))
	if e != nil {
		return 0, e
	}
	defer bl.Close()

	return bl.Insert(data)
}

/*func GetReader(site string, id int64) (io.Reader, error) {
	bl, e := bloblog.Open(filepath.Join("sites", site, "db", "files.bloblog"))
	if e != nil {
		return nil, e
	}
	defer bl.Close()

bl.Get(id)

}*/

func Get(site string, id int64) ([]byte, error) {
	bl, e := bloblog.Open(filepath.Join("sites", site, "db", "files.bloblog"))
	if e != nil {
		return nil, e
	}
	defer bl.Close()

	return bl.Get(id)
}
