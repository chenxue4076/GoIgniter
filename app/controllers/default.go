package controllers

import (
)

type MainController struct {
	BaseController
}

func (c *MainController) Index() {
	c.Data["Website"] = "localhost"
	c.Data["Email"] = "windigniter@163.com"
	c.TplName = "index.tpl"
}
