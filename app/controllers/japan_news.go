package controllers

import (
	"strconv"
	"windigniter.com/app/services"
	"github.com/astaxie/beego"
	"strings"
	"encoding/json"
	"fmt"
)

type JapanNewsController struct {
	BaseController
}

type JapanNewsDict struct {
	Def		string			`json:"def"`
	Hyouki	[]string		`json:"hyouki"`
}
type JapanNewsDictList struct {
	Dict	[]*JapanNewsDict
}

func (c *JapanNewsController) Index() {
	perPage := 2
	lang := c.CurrentLang()
	pageString := c.GetString("page")
	if pageString == "" {
		pageString = "1"
	}
	page, _ := strconv.Atoi(pageString)
	db := new(services.JapanNewsService)
	newsList, total64, err := db.JapanNewsList(perPage, page, 1)
	if err != nil {
		c.Data["Title"] = Translate(lang,"common.unknownError")
		c.Data["Content"] = Translate(lang,err.Error())
		c.Abort("Normal")
	}
	c.Data["Total"] = strconv.FormatInt(total64, 10)
	c.Data["NewsList"] = newsList
}

func (c *JapanNewsController) Show() {
	id64, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
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
	dictList := JapanNewsDictList{}
	fmt.Println("Dict is ",news.Dict)

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(news.Dict), &dat)
	if err != nil {
		fmt.Println("err is ",err)
	}
	fmt.Println(" dat is ",dat)
	fmt.Println(" result is ",dictList)
	c.Data["Dict"] = dictList
	c.Data["Title"] = news.Title
	//c.Data["Description"] = libraries.RemoveHtml(news.DescribeRuby, false)
	c.Data["Description"] = strings.Replace(beego.HTML2str(beego.Htmlunquote(news.DescribeRuby)), "\n", "", -1)
	c.Data["News"] = news
}