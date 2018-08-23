package controllers

import (
	"strconv"
	"windigniter.com/app/services"
	"fmt"
)

type BlogController struct {
	BaseController
}

// @router /blog [get]
// @router /blog/index [get]
func (c *BlogController) Index() {
	perPage := 2
	lang := c.CurrentLang()
	pageString := c.GetString("page")
	if pageString == "" {
		pageString = "1"
	}
	page, _ := strconv.Atoi(pageString)
	db := new(services.WpUsersService)
	blogs, total64, err := db.BlogList(perPage, page)
	if err != nil {
		c.Data["Title"] = Translate(lang,"common.unknownError")
		c.Data["Content"] = Translate(lang,err.Error())
		c.Abort("Normal")
	}
	dateFormat, _ := db.Options("date_format")
	permalinkStructure, _ := db.Options("permalink_structure")
	c.Data["DateFormat"] = dateFormat
	//timeTmp, _ := time.Parse("2006-01-02 15:04:05", dateFormat)
	//fmt.Println(libraries.DateFormat(timeTmp , "Y-m-d H:i:s"))
	c.Data["PermalinkStructure"] = permalinkStructure

	c.Data["Total"] = strconv.FormatInt(total64, 10)
	c.Data["Blogs"] = blogs
}

func (c *BlogController) Show() {
	id64, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	lang := c.CurrentLang()
	postName := c.Ctx.Input.Param(":postName")
	fmt.Println(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":postName"), c.Ctx.Input.Params())
	if id64 == 0 && postName == "" {
		c.Data["Title"] = Translate(lang,"common.error404")
		c.Data["Content"] = Translate(lang,"common.pageNotFound")
		c.Abort("Normal")
	}
	db := new(services.WpUsersService)
	blog, err := db.BlogDetail(id64, postName)
	if err != nil {
		c.Data["Title"] = Translate(lang,"common.unknownError")
		c.Data["Content"] = Translate(lang,err.Error())
		c.Abort("Normal")
	}
	dateFormat, _ := db.Options("date_format")
	c.Data["DateFormat"] = dateFormat

	c.Data["Title"] = blog.PostTitle
	c.Data["Author"] = blog.PostAuthor.DisplayName
	c.Data["Blog"] = blog
}