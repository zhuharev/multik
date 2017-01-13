package main

import (
	"log"
	"strings"
	"time"

	"html/template"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"pure/multik/controllers"
	"pure/multik/controllers/admin"
	"pure/multik/controllers/auth"
	"pure/multik/controllers/callback"
	"pure/multik/controllers/entity"
	"pure/multik/controllers/images"
	"pure/multik/controllers/lua"
	"pure/multik/controllers/wiki"
	"pure/multik/modules/base"
	"pure/multik/modules/middleware"
	"pure/multik/modules/setting"
)

func main() {
	controllers.GlobalInit()

	m := macaron.New()
	m.Use(func(c *macaron.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[Response time] %s %s", c.Req.URL, time.Since(start))
	})
	m.Use(macaron.Recovery())
	m.Use(Static(StaticOptions{SkipLogging: false}))
	m.Use(macaron.Renderer(macaron.RenderOptions{Layout: "layout", Funcs: []template.FuncMap{
		template.FuncMap{
			"raw":       func(in string) template.HTML { return template.HTML(in) },
			"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
			"menuEqual": func(s, prefix string) bool {
				if len(prefix) == 1 {
					return s == prefix
				} else {
					return strings.HasPrefix(s, prefix)
				}
			},
			"markdown": func(s string) template.HTML {
				return template.HTML(base.RenderMarkdownString(s))
			},
		},
	}}))

	//m.Use(macaron.Renderer(macaron.RenderOptions{Directory: "sites/*/templates"}))

	m.Use(session.Sessioner())
	m.Use(func(c *macaron.Context) {
		if strings.HasPrefix(c.Req.RequestURI, "/admin/") {
			c.SetTemplatePath("", "templates")
		} else {
			c.SetTemplatePath("", "sites/"+c.Req.Host+"/templates")
		}
	})
	m.Use(middleware.Contexter())

	m.Get(setting.Router.EntityItem, entity.Item)

	m.Group("/auth", func() {
		m.Any("/login", auth.Login)
	})

	m.Group("/admin", func() {
		m.Any("/login", admin.Login)
		m.Get("/dashboard", admin.Dashboard)

		m.Group("/", func() {
			m.Group("/eav", func() {
				m.Any("/attribute/add", admin.EavAttributeAdd)
			})
			m.Group("/rubric", func() {
				m.Get("/", admin.RubricList)
				m.Any("/add", admin.RubricCreate)
				m.Any("/:id/edit", admin.RubricEdit)
				m.Get("/:id/delete", admin.RubricDelete)
			})

			m.Group("/entity", func() {
				m.Get("/", admin.EntityList)
				m.Any("/add", admin.EntityAdd)
				m.Any("/:id/edit", admin.EntityEdit)
				m.Any("/value/add", admin.EntityValueAdd)
				m.Post("/:id/upload", admin.EntityUpload)
			})
			m.Group("/callback", func() {
				m.Get("/list", admin.CallbackFind)
			})

			m.Group("/menus", func() {
				m.Get("/", admin.MenuList)
				m.Post("/create", admin.MenuCreate)
				m.Any("/edit", admin.MenuEdit)
				m.Any("/itemcreate", admin.MenuItemCreate)
				m.Any("/itemedit", admin.MenuItemEdit)
				m.Post("/setposition", admin.MenuSetPosition)
			})

			m.Group("/widgets", func() {
				m.Any("/:slug/edit", admin.WidgetEdit)
			})

			m.Group("/files", func() {
				m.Any("/upload", admin.Upload)
				m.Get("/list", admin.FilesList)
			})

			m.Group("/posts", func() {
				m.Any("/edit", admin.PostEdit)
			})

		})

	})

	m.Group("/wiki", func() {
		m.Any("/:slug/edit", wiki.Edit)
		m.Get("/:slug/history", wiki.History)
		m.Get("/:slug", wiki.View)
	})

	m.Get("/*", controllers.Router)

	m.Any("/callback/new", callback.New)

	m.Get("/lua", lua.Lua1)

	m.Get("/", controllers.Index)

	m.Get("/img/:name", images.Show)

	m.Run(4554)
}
