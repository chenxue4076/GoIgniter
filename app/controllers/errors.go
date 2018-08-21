package controllers

type ErrorsController struct {
	BaseController
}

func (c *ErrorsController) Error404() {
	lang := c.CurrentLang()
	c.LayoutSections["HtmlHead"] = ""
	c.LayoutSections["HtmlFoot"] = ""
	c.LayoutSections["Scripts"] = ""
	c.LayoutSections["SideBar"] = ""
	if c.Data["Title"] == "" {
		c.Data["Title"] = Translate(lang, "common.pageNotFound")
	}
	if c.Data["Content"] == "" {
		c.Data["Content"] = Translate(lang, "common.pageNotFound")
	}
	c.TplName = "errors/404.html"
}

func (c *ErrorsController) ErrorNormal() {
	lang := c.CurrentLang()
	c.LayoutSections["HtmlHead"] = ""
	c.LayoutSections["HtmlFoot"] = ""
	c.LayoutSections["Scripts"] = ""
	c.LayoutSections["SideBar"] = ""
	if c.Data["Title"] == "" {
		c.Data["Title"] = Translate(lang, "common.unknownError")
	}
	if c.Data["Content"] == "" {
		c.Data["Content"] = Translate(lang, "common.unknownError")
	}
	c.TplName = "errors/Normal.html"
}