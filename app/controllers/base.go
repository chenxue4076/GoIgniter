package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) InitDefault()  {
	c.TplExt = "html"
	controller, action := c.GetControllerAndAction()
	c.TplName = strings.Replace(strings.ToLower(controller), "controller", "", -1) + "/" + strings.ToLower(action) + "." +  c.TplExt
	c.Layout = "layout/common."+c.TplExt

	c.Data["title"] = ""
}

