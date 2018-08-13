package controllers

import (
		"html/template"
	)

type UserController struct {
	BaseController
	//InitDefault( *BaseController)
}

func (c *UserController) Prepare() {
	c.InitDefault()
}

func (c *UserController) LoginForm() {
	c.Data["XsrfData"] = template.HTML(c.XSRFFormHTML())
	//c.TplExt = "html"
}

func (c *UserController) LoginPost() {
	c.TplName = "user/login"
}