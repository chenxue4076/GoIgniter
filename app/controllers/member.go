package controllers

import "fmt"

type MemberController struct {
	BaseController
}

func (c *MemberController) Center() {
	username := c.Ctx.Input.Param(":username")
	fmt.Println(username)
}