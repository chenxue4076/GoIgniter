package controllers

import (
	"net/http"
	"fmt"
)

type UserController struct {
	BaseController
}

func (c *UserController) Login() {
	if c.Ctx.Request.Method == http.MethodPost {	//POST Login deal
		//a := c.GetString("username")
		//b := c.Input().Get("password")
		//fmt.Println("username:", a, " password:", b)
		username := c.Input().Get("username")
		password := c.Input().Get("password")
		fmt.Println("username:", username, " password:", password)
		//return
	}
	//c.LayoutSections["HeaderMeta"] = "user/headermeta.html"
	//c.LayoutSections["HtmlHead"] = ""
	c.LayoutSections["HtmlFoot"] = ""
	//c.LayoutSections["Scripts"] = "user/scripts.html"
}

/*func (c *UserController) LoginPost() {
	c.TplName = "user/login"
}*/