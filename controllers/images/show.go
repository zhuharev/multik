package images

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"github.com/Unknwon/com"
	//"pure/multik/models"
	"github.com/disintegration/imaging"
	"github.com/sisteamnik/guseful/md5"
	"pure/multik/modules/filedb"
	"pure/multik/modules/middleware"
)

const (
	WidthLimit  = 3000
	HeightLimit = 3000
)

func Show(c *middleware.Context) {
	var (
		path    = c.Req.Request.URL.Path
		resized bool
		id      int64

		width  = 100
		height = 100
	)

	path = strings.TrimSuffix(path, ".jpg")

	_, fileName := filepath.Split(path)
	if strings.Contains(fileName, "_") {
		resized = true

		arr := strings.Split(fileName, "_")
		id = com.StrTo(arr[0]).MustInt64()
		sizes := strings.Split(arr[1], "x")

		width = com.StrTo(sizes[0]).MustInt()
		if width > WidthLimit {
			width = WidthLimit
		}
		height = com.StrTo(sizes[1]).MustInt()
		if height > HeightLimit {
			height = HeightLimit
		}
	} else {
		id = com.StrTo(fileName).MustInt64()
	}

	fName := fmt.Sprintf("%d_%dx%d", id, width, height)
	hash := md5.Hash(fName)
	fPath := filepath.Join("sites", c.E.Site(), "public", "img", fmt.Sprintf("%s/%s/%s", string(hash[0]), string(hash[1]), string(hash[2])))
	if !com.IsExist(fPath) {
		e := os.MkdirAll(fPath, 0777)
		if e != nil {
			panic(e)
		}
	}

	fullName := filepath.Join(fPath, fName)

	if com.IsExist(fullName) {
		f, e := os.Open(fullName)
		if e != nil {
			panic(e)
		}
		defer f.Close()
		c.Resp.Header().Set("Content-Type", "image/jpeg")
		c.Resp.WriteHeader(200)

		_, e = io.Copy(c.Resp, f)
		if e != nil {
			panic(e)
		}
		return
	}

	file, e := c.E.FilesGet(id)
	if e != nil {
		panic(e)
	}

	bts, e := filedb.Get(c.E.Site(), file.FileId)
	if e != nil {
		panic(e)
	}

	if resized {
		rdr := bytes.NewReader(bts)

		img, _, e := image.Decode(rdr)
		if e != nil {
			panic(e)
		}

		img = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

		// img = resize.Thumbnail(width, height, img, resize.Bilinear)

		buf := bytes.NewBuffer(nil)

		e = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
		if e != nil {
			panic(e)
		}

		bts = buf.Bytes()
	}

	e = com.WriteFile(fullName, bts)
	if e != nil {
		panic(e)
	}

	c.Resp.Header().Set("Content-Type", "image/jpeg")
	c.Resp.WriteHeader(200)
	c.Resp.Write(bts)
}
