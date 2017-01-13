package lua

import (
	"github.com/yuin/gopher-lua"
	"pure/multik/modules/middleware"
)

func Lua1(c *middleware.Context) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}

	L = lua.NewState()
	defer L.Close()
	if err := L.DoFile("sites/" + c.Req.Host + "/lua/hello.lua"); err != nil {
		panic(err)
	}
}
