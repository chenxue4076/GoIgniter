package controllers

import (
	"strconv"
	"windigniter.com/app/services"
	"github.com/astaxie/beego"
	"strings"
	"encoding/json"
	"fmt"
	"windigniter.com/app/libraries"
	"net/http"
	"time"
	"github.com/astaxie/beego/validation"
	"io/ioutil"
	"bytes"
)

type JapanNewsController struct {
	BaseController
}
// japan news list
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
// japan detail page
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
// struct japan news
type JapanNewsTopListItem struct {
	TopPriorityNumber					string			`json:"top_priority_number"`
	TopDisplayFlag						bool			`json:"top_display_flag"`
	NewsId								string			`json:"news_id"`
	NewsPrearrangedTime					time.Time		`json:"news_prearranged_time;type(datetime)"`
	Title								string			`json:"title"`
	TitleWithRuby						string			`json:"title_with_ruby"`
	OutlineWithRuby						string			`json:"outline_with_ruby"`
	NewsFileVer							bool			`json:"news_file_ver"`
	NewsPublicationStatus				bool			`json:"news_publication_status"`
	HasNewsWebImage						bool			`json:"has_news_web_image"`
	HasNewsWebMovie						bool			`json:"has_news_web_movie"`
	HasNewsEasyImage					bool			`json:"has_news_easy_image"`
	HasNewsEasyMovie					bool			`json:"has_news_easy_movie"`
	HasNewsEasyVoice					bool			`json:"has_news_easy_voice"`
	NewsWebImageUri						string			`json:"news_web_image_uri"`
	NewsWebMovieUri						string			`json:"news_web_movie_uri"`
	NewsEasyImageUri					string			`json:"news_easy_image_uri"`
	NewsEasyMovieUri					string			`json:"news_easy_movie_uri"`
	NewsEasyVoiceUri					string			`json:"news_easy_voice_uri"`
}

// background crawl easy news
//http://www3.nhk.or.jp/news/easy/top-list.json?_=1484116080539 //最新7条
//http://www3.nhk.or.jp/news/easy/news-list.json?_=1484116080540    一周列表
//http://www3.nhk.or.jp/news/easy/k10010833901000/k10010833901000.out.dic?date=1484119973650    //字典字段
func (c *JapanNewsController) Crawl() {
	url := "http://www3.nhk.or.jp/news/easy/top-list.json" 	//?_='. substr($microtime, -10) . substr($microtime, 2, 3);
	microTime := strconv.FormatInt(time.Now().UnixNano(),  10)
	valid := validation.Validation{}
	resp, err := http.Get(url + "?_="+ beego.Substr(microTime, 0, 13))
	if err != nil {
		fmt.Println(err)
		valid.SetError("", err.Error())
		c.Data["Error"] = valid.Errors
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		valid.SetError("", err.Error())
		c.Data["Error"] = valid.Errors
	}
	body = bytes.TrimPrefix(body, []byte{239,187,191})	//去掉网页头部BOM
	topList := []JapanNewsTopListItem{}
	err = json.Unmarshal(body, &topList)
	if err != nil {
		fmt.Println(err)
		valid.SetError("", err.Error())
		c.Data["Error"] = valid.Errors
	}
	fmt.Println(topList)
}

