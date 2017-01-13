package middleware

import (
	"bytes"
	"html/template"
	"io/ioutil"
	//"reflect"

	//"github.com/fatih/color"
	anko_core "github.com/mattn/anko/builtins"
	"github.com/mattn/anko/vm"
)

type Widget struct {
	r func(string) template.HTML
}

func (w Widget) Render(name string) template.HTML {
	if w.r == nil {
		return ""
	}
	return w.r(name)
}

func makeWidget(c *Context) Widget {

	f := func(name string) template.HTML {
		c := c

		basePath := "sites/" + c.Req.Host + "/widgets/" + name + "/"

		t, err := template.ParseFiles(basePath+"template.tmpl", basePath+"error.tmpl")
		if err != nil {
			panic(err)
		}
		val := template.HTML("")
		if val, err = DoFile(c, name, basePath+"widget.anko", t); err != nil {
			panic(err)
		}

		//color.Green("%s", val)
		return val

		//wr := bytes.NewBuffer(nil)

		//err = t.ExecuteTemplate(wr, "template.tmpl", val)
		//if err != nil {
		//	panic(err)
		//}

		//return template.HTML(wr.String())

	}

	wid := Widget{
		r: f,
	}
	return wid
}

func DoFile(c *Context, name, path string, t *template.Template) (template.HTML, error) {

	var (
		result template.HTML
	)

	bts, e := ioutil.ReadFile(path)
	if e != nil {
		return result, e
	}

	env := vm.NewEnv()
	anko_core.LoadAllBuiltins(env)

	key := "widget." + name + "." + "data"

	data, e := c.E.Get(key)
	if e != nil {
		return result, e
	}

	e = env.DefineGlobal("widgetData", string(data))
	if e != nil {
		return result, e
	}

	e = env.DefineGlobal("render", func(data interface{}) template.HTML {
		wr := bytes.NewBuffer(nil)
		e = t.ExecuteTemplate(wr, "template.tmpl", data)
		if e != nil {
			return template.HTML(e.Error())
		}
		return template.HTML(wr.String())
	})

	e = env.DefineGlobal("widgetName", name)
	if e != nil {
		return result, e
	}

	e = env.DefineGlobal("context", c)
	if e != nil {
		return result, e
	}

	val, e := env.Execute(string(bts))
	if e != nil {
		return result, e
	}
	return template.HTML(val.String()), nil

	//val, e := env.Get("result")
	//return val, e
}
