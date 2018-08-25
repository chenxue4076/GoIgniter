package controllers

import (
	"strconv"
	"windigniter.com/app/services"
	"github.com/astaxie/beego"
	"strings"
	"encoding/json"
	"fmt"
	"windigniter.com/app/libraries"
)

type JapanNewsController struct {
	BaseController
}

type JapanNewsDict struct {
	Def		string			`json:"def"`
	Hyouki	[]string		`json:"hyouki"`
}

func (c *JapanNewsController) Index() {
	perPage := 5
	lang := c.CurrentLang()
	pageString := c.GetString("page")
	if pageString == "" {
		pageString = "1"
	}
	page, _ := strconv.Atoi(pageString)
	db := new(services.JapanNewsService)
	newsList, total, err := db.JapanNewsList(perPage, page, 1)
	if err != nil {
		c.Data["Title"] = Translate(lang,"common.unknownError")
		c.Data["Content"] = Translate(lang,err.Error())
		c.Abort("Normal")
	}
	c.Data["Total"] = total
	c.Data["NewsList"] = newsList
	c.Data["Title"] = Translate(lang,"japannews.title")
	//fmt.Println(c.Ctx.Request.URL, c.Ctx.Request.RequestURI, c.Ctx.Request.URL.Path)
	c.Data["Pagination"] = libraries.PageList(total, page, perPage, c.Ctx.Request.URL.Path, 3, lang)
}

func (c *JapanNewsController) Show() {
	var id64 int64
	id := c.GetString("id")
	if id == "" {
		id64, _ = strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	} else {
		id64, _ = strconv.ParseInt(id, 10, 64)
	}
	lang := c.CurrentLang()
	if id64 == 0 {
		c.Data["Title"] = Translate(lang,"common.error404")
		c.Data["Content"] = Translate(lang,"common.pageNotFound")
		c.Abort("Normal")
	}
	db := new(services.JapanNewsService)
	news, err := db.JapanNewsDetail(id64)
	if err != nil {
		c.Data["Title"] = Translate(lang,"common.unknownError")
		c.Data["Content"] = Translate(lang,err.Error())
		c.Abort("Normal")
	}

	//dict json decode

	var dictMap map[string]interface{}
	err = json.Unmarshal([]byte(news.Dict), &dictMap)
	if err != nil {
		fmt.Println("err is ",err)
	}
	c.Data["Dict"] = dictMap
	c.Data["Title"] = news.Title
	//c.Data["Description"] = libraries.RemoveHtml(news.DescribeRuby, false)
	c.Data["Description"] = strings.Replace(beego.HTML2str(beego.Htmlunquote(news.DescribeRuby)), "\n", "", -1)
	c.Data["News"] = news
}