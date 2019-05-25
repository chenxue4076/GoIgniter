package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"windigniter.com/app/libraries"
	"windigniter.com/app/models"
	"windigniter.com/app/services"
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
		c.Data["Title"] = Translate(lang, "common.unknownError")
		c.Data["Content"] = Translate(lang, err.Error())
		c.Abort("Normal")
	}
	c.Data["Total"] = total
	c.Data["NewsList"] = newsList
	c.Data["Title"] = Translate(lang, "japannews.title")
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
		c.Data["Title"] = Translate(lang, "common.error404")
		c.Data["Content"] = Translate(lang, "common.pageNotFound")
		c.Abort("Normal")
	}
	db := new(services.JapanNewsService)
	news, err := db.JapanNewsDetail(id64)
	if err != nil {
		c.Data["Title"] = Translate(lang, "common.unknownError")
		c.Data["Content"] = Translate(lang, err.Error())
		c.Abort("Normal")
	}

	//dict json decode

	var dictMap map[string]interface{}
	err = json.Unmarshal([]byte(news.Dict), &dictMap)
	if err != nil {
		fmt.Println("err is ", err)
	}
	c.Data["Dict"] = dictMap
	c.Data["Title"] = news.Title
	//c.Data["Description"] = libraries.RemoveHtml(news.DescribeRuby, false)
	c.Data["Description"] = strings.Replace(beego.HTML2str(beego.Htmlunquote(news.DescribeRuby)), "\n", "", -1)
	c.Data["News"] = news
}

//japan news radio player
func (c *JapanNewsController) Player() {
	c.LayoutSections["HeaderMeta"] = ""
	c.LayoutSections["HtmlHead"] = ""
	c.LayoutSections["HtmlFoot"] = ""
	c.LayoutSections["Scripts"] = ""
	c.LayoutSections["SideBar"] = ""
}

// struct japan news
type JapanNewsTopListItem struct {
	TopPriorityNumber     string    `json:"top_priority_number"`
	TopDisplayFlag        bool      `json:"top_display_flag"`
	NewsId                string    `json:"news_id"`
	NewsPrearrangedTime   time.Time `json:"news_prearranged_time;type(datetime)"`
	Title                 string    `json:"title"`
	TitleWithRuby         string    `json:"title_with_ruby"`
	OutlineWithRuby       string    `json:"outline_with_ruby"`
	NewsFileVer           bool      `json:"news_file_ver"`
	NewsPublicationStatus bool      `json:"news_publication_status"`
	HasNewsWebImage       bool      `json:"has_news_web_image"`
	HasNewsWebMovie       bool      `json:"has_news_web_movie"`
	HasNewsEasyImage      bool      `json:"has_news_easy_image"`
	HasNewsEasyMovie      bool      `json:"has_news_easy_movie"`
	HasNewsEasyVoice      bool      `json:"has_news_easy_voice"`
	NewsWebImageUri       string    `json:"news_web_image_uri"`
	NewsWebMovieUri       string    `json:"news_web_movie_uri"`
	NewsEasyImageUri      string    `json:"news_easy_image_uri"`
	NewsEasyMovieUri      string    `json:"news_easy_movie_uri"`
	NewsEasyVoiceUri      string    `json:"news_easy_voice_uri"`
}

// background crawl easy news
//http://www3.nhk.or.jp/news/easy/top-list.json?_=1484116080539 //最新7条
//http://www3.nhk.or.jp/news/easy/news-list.json?_=1484116080540    一周列表
//http://www3.nhk.or.jp/news/easy/k10010833901000/k10010833901000.out.dic?date=1484119973650    //字典字段
func (c *JapanNewsController) Crawl() {
	url := "http://www3.nhk.or.jp/news/easy/top-list.json" //?_='. substr($microtime, -10) . substr($microtime, 2, 3);
	microTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	valid := validation.Validation{}
	resp, err := http.Get(url + "?_=" + beego.Substr(microTime, 0, 13))
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
	body = bytes.TrimPrefix(body, []byte{239, 187, 191}) //去掉网页头部BOM
	topList := []JapanNewsTopListItem{}
	err = json.Unmarshal(body, &topList)
	if err != nil {
		fmt.Println(err)
		valid.SetError("", err.Error())
		c.Data["Error"] = valid.Errors
	}
	//fmt.Println(topList)
	db := new(services.JapanNewsService)
	for _, item := range topList {
		if item.NewsId == "" {
			continue
		}
		hasError := false
		easyNews, err := db.JapanEasyNewsDetail(item.NewsId)
		if err != nil {
			if err.Error() == "common.ormErrNoRows" {
				easyNews = models.JapanEasyNews{NewsId: item.NewsId, NewsPrearrangedTime: item.NewsPrearrangedTime, Title: item.Title, TitleWithRuby: item.TitleWithRuby, OutlineWithRuby: item.OutlineWithRuby, NewsWebImageUri: item.NewsWebImageUri, NewsWebMovieUri: item.NewsWebMovieUri, NewsEasyVoiceUri: item.NewsEasyVoiceUri, Status: 0}
				_, err := db.SaveEasyNews(easyNews)
				if err != nil {
					hasError = true
					fmt.Println("save japan easy news err :", err)
				}
			} else {
				hasError = true
				fmt.Println("for topList hasEasy err :", err)
			}
		}
		if !hasError && easyNews.Status == 0 { //has not saved to japan news
			content, err := crawlJapanNewsContent(item.NewsId)
			if err != nil {
				fmt.Println("crawlJapanNewsContent returns err", err)
				continue
			}
			dictJson, err := crawlJapanNewsDict(item.NewsId)
			if err != nil {
				fmt.Println("crawlJapanNewsDict returns err", err)
				continue
			}
			japanNews := models.JapanNews{AuthorId: 1, NewsId: easyNews.NewsId, Title: easyNews.Title, TitleRuby: easyNews.TitleWithRuby, DescribeRuby: easyNews.OutlineWithRuby, Views: 0, Featured: easyNews.NewsWebImageUri, Media: easyNews.NewsEasyVoiceUri, Content: content, Dict: dictJson, Ding: 1, Cai: 1, Pubdate: time.Now(), Status: 1}
			//fmt.Println(japanNews)
			id, err := db.SaveJapanNews(japanNews)
			if err != nil {
				fmt.Println("save japan news err", err)
				continue
			}
			if id > 0 {
				easyNews.Status = 1
				err = db.UpdateEasyNews(easyNews, "status")
				if err != nil {
					fmt.Println("update japan easy news status err", easyNews.NewsId, err)
					continue
				}
			}
		}
	}
}

// view japan news content in japan news web site
func (c *JapanNewsController) NewsContent() {
	valid := validation.Validation{}
	newsId := c.GetString("newsid")
	if newsId == "" {
		valid.SetError("", "No newsId")
		c.Data["Error"] = valid.Errors
	}
	/*result, err := crawlJapanNewsContent(newsId)
	if err != nil {
		valid.SetError("", err.Error())
		c.Data["Error"] = valid.Errors
	}
	c.Data["Content"] = result*/
	result, err := crawlJapanNewsDict(newsId)
	if err != nil {
		valid.SetError("", err.Error())
		c.Data["Error"] = valid.Errors
	}
	c.Data["Content"] = result
}

// analyze japan news content
func crawlJapanNewsContent(newsId string) (result string, err error) {
	preNewsUrl := "https://www3.nhk.or.jp/news/easy/" + newsId + "/" + newsId
	articleUrl := preNewsUrl + ".html"
	fmt.Println("content url ", articleUrl)
	respArticle, err := http.Get(articleUrl)
	if err != nil {
		fmt.Println("get article resp error ", err)
		return result, err
	}
	defer respArticle.Body.Close()
	// japan news content
	doc, err := html.Parse(respArticle.Body)
	if err != nil {
		fmt.Println("get body article error ", err)
		return result, err
	}
	hasFind := false
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			//fmt.Println(n)
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "js-article-body" {
					hasFind = true
					buf := new(bytes.Buffer)
					err = html.Render(buf, n)
					if err != nil {
						fmt.Println("rand err", err)
					}
					result = buf.String()
					break
				}
			}
		}
		if !hasFind {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)
	if hasFind {
		return result, nil
	}
	return result, err
}

// get japan news dictionary
func crawlJapanNewsDict(newsId string) (result string, err error) {
	microTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	preNewsUrl := "https://www3.nhk.or.jp/news/easy/" + newsId + "/" + newsId
	dictUrl := preNewsUrl + ".out.dic?date=" + beego.Substr(microTime, 0, 13)
	fmt.Println(dictUrl)
	respDict, err := http.Get(dictUrl)
	if err != nil {
		fmt.Println("get article resp error ", err)
		return result, err
	}
	defer respDict.Body.Close()
	body, err := ioutil.ReadAll(respDict.Body)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	bodyString := strings.Replace(string(body), "\n", "", -1)
	bodyString = strings.Replace(bodyString, " ", "", -1)
	jsonByte := []byte(bodyString)
	//fmt.Println(bodyString)
	//fmt.Println([]byte(bodyString))
	jsonData := jsonByte[21:(len(jsonByte) - 2)]
	//fmt.Println(jsonData)
	return string(jsonData), err
}
