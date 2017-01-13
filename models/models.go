package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/fatih/color"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhuharev/eav"
)

var (
	xs = map[string]*Engine{}
	re = regexp.MustCompile(`sites/(.*\..*)/db`)

	tables = []interface{}{}
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

func GetEngine(host string) *Engine {
	return xs[host]
}

func NewContext() {

	tables = append(tables,
		new(Entity),
		new(eav.Value),
		new(eav.Attribute),
		new(Post),
		new(Callback),
		new(Rubric),
		new(RubricItem),
		new(Menu),
		new(MenuItem),
		new(KV),
		new(File),
		new(Route),
		new(WikiPage),
		new(WikiDelta))

	matches, e := filepath.Glob("sites/*/db")
	if e != nil {
		log.Fatalln(e)
	}
	for _, v := range matches {
		color.Green("Load %s", v)
		arr := re.FindStringSubmatch(v)
		if len(arr) != 2 {
			//color.Red("%s", )
			log.Fatalln("its bad, arr len must be 2")
		}
		host := arr[1]
		color.Green("%s host found", host)
		eng, e := xorm.NewEngine("sqlite3", filepath.Join(v, "db.sqlite"))
		if e != nil {
			log.Fatalln(e)
		}
		eng.Sync2(tables...)

		f, e := os.OpenFile(filepath.Join(v, "db.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if e != nil {
			log.Fatalln(e)
		}
		logger := xorm.NewSimpleLogger(f)
		eng.SetLogger(logger)
		eng.ShowDebug = true
		eng.ShowSQL = true
		eng.ShowErr = true
		eng.ShowInfo = true
		eng.ShowWarn = true

		en := new(Engine)
		en.Engine = eng
		en.site = host
		xs[host] = en
	}
}

func (x *Engine) Save(bean interface{}) (e error) {
	val := reflect.ValueOf(bean)
	val = val.Elem()
	idVal := val.FieldByName("Id")
	if !idVal.IsValid() {
		return fmt.Errorf("id value is nil")
	}
	id := idVal.Interface().(int64)

	if id == 0 {
		_, e = x.InsertOne(bean)
		if e != nil {
			return
		}
		//e = SaveTags(bean)
		return
	} else {
		_, e = x.Id(id).Update(bean)
		if e != nil {
			return
		}
		//e = SaveTags(bean)
		return
	}
	return
}
