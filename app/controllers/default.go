package controllers

import (
	"net/http"
)

type MainController struct {
	BaseController
}

func (c *MainController) Index() {
	//http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, "/japan", 302)
	//http.RedirectHandler("/japan", 301)
	//c.Data["Website"] = "localhost"
	//c.Data["Email"] = "windigniter@163.com"
	c.TplName = "index.tpl"
}
