package controllers

import "github.com/astaxie/beego"

type MemberController struct {
	BaseController
}

func (c *MemberController) Center() {
	lang := c.CurrentLang()
	c.Data["Title"] = Translate(lang, "member.userCenter")
	c.Data["Refer"] = beego.URLFor("UserController.Login")
}

func (c *MemberController) Index() {

}
