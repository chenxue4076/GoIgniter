package controllers

import "github.com/astaxie/beego"

type MemberController struct {
	BaseController
}

func (c *MemberController) Center() {
	c.Data["Refer"] = beego.URLFor("UserController.Login")


}

func (c *MemberController) Index()  {
	
}