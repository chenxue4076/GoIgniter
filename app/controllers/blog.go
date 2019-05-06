package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"windigniter.com/app/libraries"
	"windigniter.com/app/models"
	"windigniter.com/app/services"
)

type BlogController struct {
	BaseController
}
/**	sidebar info, call can call
 */
func sideBar( c *BlogController) {
	db := new(services.WpUsersService)
	tagList, _ := db.TagList("")
	if tagList != nil {
		c.Data["TagList"] = tagList
	}
	archiveList, _ := db.ArchiveList()
	if archiveList != nil {
		c.Data["ArchiveList"] = archiveList
	}
	categoryList, _ := db.TagList("category")
	if categoryList != nil {
		c.Data["CategoryList"] = categoryList
	}

}

// @router /blog [get]
// @router /blog/index [get]
func (c *BlogController) Index() {
	lang := c.CurrentLang()
	db := new(services.WpUsersService)
	perPage := 5
	pageString := c.GetString("page")
	if pageString == "" {
		pageString = "1"
	}
	page, _ := strconv.Atoi(pageString)
	slug := c.Ctx.Input.Param(":slug")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")

	var blogs []*models.WpPosts
	var total int
	var err error

	if slug != "" {
		tag, err := db.TagInfo(slug)
		if err != nil {
			c.Data["Title"] = Translate(lang, "common.unknownError")
			c.Data["Content"] = Translate(lang, err.Error())
			c.Abort("Normal")
		}
		blogs, total, err = db.BlogListByTermId(tag.TermId, perPage, page)
		if err != nil {
			c.Data["Title"] = Translate(lang, "common.unknownError")
			c.Data["Content"] = Translate(lang, err.Error())
			c.Abort("Normal")
		}

	} else if year != "" && month != "" {
		blogs, total, err = db.BlogListByDate(year, month, perPage, page)
		if err != nil {
			c.Data["Title"] = Translate(lang, "common.unknownError")
			c.Data["Content"] = Translate(lang, err.Error())
			c.Abort("Normal")
		}
	} else {
		blogs, total, err = db.BlogList(perPage, page)
		if err != nil {
			c.Data["Title"] = Translate(lang, "common.unknownError")
			c.Data["Content"] = Translate(lang, err.Error())
			c.Abort("Normal")
		}
	}
	dateFormat, _ := db.Options("date_format")
	permalinkStructure, _ := db.Options("permalink_structure")
	c.Data["DateFormat"] = dateFormat
	//timeTmp, _ := time.Parse("2006-01-02 15:04:05", dateFormat)
	//fmt.Println(libraries.DateFormat(timeTmp , "Y-m-d H:i:s"))
	c.Data["PermalinkStructure"] = permalinkStructure

	//侧边栏
	sideBar(c)
	c.Data["Total"] = total
	c.Data["Blogs"] = blogs
	c.Data["Title"] = Translate(lang, "common.Blog")
	c.Data["Pagination"] = libraries.PageList(total, page, perPage, c.Ctx.Request.URL.Path, 3, lang)
}

func (c *BlogController) Show() {
	lang := c.CurrentLang()
	id64, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	postName := c.Ctx.Input.Param(":postName")
	fmt.Println(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":postName"), c.Ctx.Input.Params())
	if id64 == 0 && postName == "" {
		c.Data["Title"] = Translate(lang, "common.error404")
		c.Data["Content"] = Translate(lang, "common.pageNotFound")
		c.Abort("Normal")
	}
	db := new(services.WpUsersService)
	blog, err := db.BlogDetail(id64, postName)
	if err != nil {
		c.Data["Title"] = Translate(lang, "common.unknownError")
		c.Data["Content"] = Translate(lang, err.Error())
		c.Abort("Normal")
	}
	tags, _ := db.Tags(blog.Id, "")
	if tags != nil {
		c.Data["Tags"] = tags
		var keywords []string
		for _, tag := range tags {
			keywords = append(keywords, tag.Name)
		}
		c.Data["Keywords"] = strings.Join(keywords, ",")
	}
	dateFormat, _ := db.Options("date_format")
	c.Data["DateFormat"] = dateFormat

	//侧边栏
	sideBar(c)
	c.Data["Title"] = blog.PostTitle
	c.Data["Author"] = blog.PostAuthor.DisplayName
	c.Data["Description"] = blog.PostExcerpt
	c.Data["Blog"] = blog
}

func (c *BlogController) Tag() {
	lang := c.CurrentLang()
	slug := c.Ctx.Input.Param(":slug")
	if slug == "" {
		c.Data["Title"] = Translate(lang, "common.error404")
		c.Data["Content"] = Translate(lang, "common.pageNotFound")
		c.Abort("Normal")
	}
	db := new(services.WpUsersService)
	tag, err := db.TagInfo(slug)
	if err != nil {
		c.Data["Title"] = Translate(lang, "common.unknownError")
		c.Data["Content"] = Translate(lang, err.Error())
		c.Abort("Normal")
	}
	perPage := 5
	pageString := c.GetString("page")
	if pageString == "" {
		pageString = "1"
	}
	page, _ := strconv.Atoi(pageString)
	blogs, total, err := db.BlogListByTermId(tag.TermId, perPage, page)
	if err != nil {
		c.Data["Title"] = Translate(lang, "common.unknownError")
		c.Data["Content"] = Translate(lang, err.Error())
		c.Abort("Normal")
	}
	dateFormat, _ := db.Options("date_format")
	permalinkStructure, _ := db.Options("permalink_structure")
	c.Data["DateFormat"] = dateFormat
	//timeTmp, _ := time.Parse("2006-01-02 15:04:05", dateFormat)
	//fmt.Println(libraries.DateFormat(timeTmp , "Y-m-d H:i:s"))
	c.Data["PermalinkStructure"] = permalinkStructure

	c.TplName = strings.Replace(c.TplName, "tag", "index", -1)

	//侧边栏
	sideBar(c)
	c.Data["Total"] = total
	c.Data["Blogs"] = blogs
	c.Data["Tag"] = tag
	c.Data["Title"] = tag.Name
	c.Data["Pagination"] = libraries.PageList(total, page, perPage, c.Ctx.Request.URL.Path, 3, lang)
}
