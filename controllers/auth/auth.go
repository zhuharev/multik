package auth

import (
	//"pure/multik/models"
	"pure/multik/modules/middleware"
)

func Login(c *middleware.Context) {
	c.HTML(200, "auth/login")
}

func Logout(c *middleware.Context) {

}
