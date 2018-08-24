package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"github.com/beego/i18n"
	"sync"
	"net/http"
	"html/template"
	"time"
	"windigniter.com/app/libraries"
)

var once sync.Once

type BaseController struct {
	beego.Controller
}

type JsonOut struct {
	Err bool
	Result JsonMessage
	Redirect string
}
type JsonMessage struct {
	Message, Key string
}

type SessionUser struct {
	Uid				int64
	UserLogin		string
	UserNicename	string
	UserEmail		string
	UserRegistered	time.Time
	DisplayName		string
}


//var langs = []string {"zh-CN", "en-US"}
//var langs []string

// all controllers init
func (c *BaseController) Prepare()  {
	//multi language load
	once.Do(libraries.LoadLangs)

	//page params
	c.TplExt = "html"
	controllerPre, actionPre := c.GetControllerAndAction()
	controller := strings.Replace(strings.ToLower(controllerPre), "controller", "", -1)
	action := strings.ToLower(actionPre)

	//page frame shows
	c.TplName = controller + "/" + action + "." +  c.TplExt
	c.Layout = "layout/common."+c.TplExt
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HeaderMeta"] = controller + "/layout_meta."+c.TplExt
	c.LayoutSections["HtmlHead"] = controller + "/layout_header."+c.TplExt
	c.LayoutSections["HtmlFoot"] = controller + "/layout_footer."+c.TplExt
	c.LayoutSections["Scripts"] = controller + "/layout_scripts."+c.TplExt
	c.LayoutSections["SideBar"] = ""

	//page data
	c.Data["Lang"] = c.CurrentLang()

	//XSRF attack protect
	c.Data["XsrfData"] = template.HTML(c.XSRFFormHTML())
	if c.Ctx.Request.Method == http.MethodPost {	//XSRF filter
		c.CheckXSRFCookie()
	}
	//active menu
	c.Data["ActiveClass"] = controller
	c.Data["Refer"] = c.Ctx.Request.RequestURI
	//if login
	c.Data["CurrentUser"] = c.GetSession("userInfo")
	//runmode
	c.Data["RunMode"] = beego.AppConfig.String("RunMode")
}

func (c *BaseController) CurrentLang() string {
	hasCookie := false
	// 1. Check URL arguments.
	lang := c.Input().Get("lang")
	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = c.Ctx.GetCookie("lang")
		hasCookie = true
	}
	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		hasCookie = false
	}
	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := c.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}
	// 4. Default language is English.
	if len(lang) == 0 {
		langs, _ := libraries.CurrentDirs("resources/lang")
		lang = langs[0]
	}
	if !hasCookie {
		c.Ctx.SetCookie("lang", lang, 1<<31-1, "/")
	}
	return lang
}

func Translate(lang,input string, args ...interface{}) string {
	return i18n.Tr(lang, input, args)
}

