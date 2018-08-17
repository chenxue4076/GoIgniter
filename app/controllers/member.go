package controllers

import "fmt"

type MemberController struct {
	BaseController
}

func (c *MemberController) Center() {
	username := c.GetSession("username")
	uid := c.GetSession("uid")
	fmt.Println(username, uid)
}

func (c *MemberController) Index()  {
	
}