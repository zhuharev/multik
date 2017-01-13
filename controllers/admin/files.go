package admin

import (
	"bytes"
	//"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "image/gif"
	_ "image/png"

	"github.com/fatih/color"
	"pure/multik/models"
	//"pure/multik/modules/filedb"
	"pure/multik/modules/middleware"
)

func upload(c *middleware.Context) (id int64, e error) {

	var (
		file = new(models.File)
	)

	c.Req.ParseMultipartForm(32 << 20)
	fileRdr, _, e := c.Req.FormFile("file")
	if e != nil {
		return
	}
	defer fileRdr.Close()

	bts, e := ioutil.ReadAll(fileRdr)
	if e != nil {
		return
	}
	rdr := bytes.NewReader(bts)

	img, _, e := image.Decode(rdr)
	if e != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	e = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	if e != nil {
		return
	}

	e = c.E.Upload(file, buf.Bytes())
	if e != nil {
		return
	}

	/*c.JSON(200, map[string]interface{}{
		"success":  true,
		"err":      e,
		"file":     file,
		"filelink": fmt.Sprintf("/img/%d_1024x512.jpg", file.Id),
	})*/
	return file.Id, nil
}

func Upload(c *middleware.Context) {

	if c.Req.Method == "GET" {
		c.HTML(200, "admin/files/upload")
		return
	}

	id, e := upload(c)
	if e != nil {
		color.Red("%s", e.Error())
	}

	c.JSON(200, id)
}

func FilesList(c *middleware.Context) {

	items, e := c.E.FilesList()
	if e != nil {
		color.Red("%s", e.Error())
	}

	c.Data["items"] = items
	c.HTML(200, "admin/files/list")
}
