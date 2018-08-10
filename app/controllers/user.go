package controllers

import (
	"github.com/astaxie/beego"
	"html/template"
)

type UserController struct {
	beego.Controller
}

func initDefault(c *UserController)  {
	c.TplExt = "html"
}

func (c *UserController) LoginForm() {
	initDefault(c)
	c.Data["XsrfData"] = template.HTML(c.XSRFFormHTML())

	c.TplName = "user/login"
}

func (c *UserController) LoginPost() {
	c.TplName = "user/login"
}