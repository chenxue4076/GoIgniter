package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"strings"
	"io/ioutil"
	"path/filepath"
	"github.com/beego/i18n"
	"sync"
	"net/http"
	"html/template"
	"time"
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


var langs = []string {"zh-CN", "en-US"}

// all controllers init
func (c *BaseController) Prepare()  {
	//multi language load
	once.Do(loadLangs)

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

	c.Data["Refer"] = c.Ctx.Request.RequestURI
	//TODO	 //unicode.In(action, []string{"login"})
	/*if controller == "user" &&  action == "login" {
		c.Data["Refer"] = ""
	}*/


	//if login
	c.Data["User"] = c.GetSession("userInfo")
	//if c.GetSession("userInfo") != nil {
		//c.Data["User"] = SessionUser{c.GetSession("Uid").(int64), c.GetSession("UserLogin").(string), c.GetSession("UserNicename").(string), c.GetSession("UserEmail").(string), c.GetSession("UserRegistered").(string), c.GetSession("DisplayName").(string)}
	// }
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
		lang = langs[0]
	}
	if !hasCookie {
		c.Ctx.SetCookie("lang", lang, 1<<31-1, "/")
	}
	return lang
}

func Translate(lang,input string) string {
	return i18n.Tr(lang, input)
}

func loadLangs()  {
	//language choose
	//langs := []string {"zh-CN", "en-US"}
	for _, lang := range langs {
		langData := make([]byte, 0)
		//beego.Trace("Loading language: " + lang)
		filepath.Walk("resources/lang/"+lang, func(path string, f os.FileInfo, err error) error {
			//fmt.Println(path, f.Name(), err)
			if f != nil && ! f.IsDir() {
				tempData, e := ReadFile(path)
				if e != nil {
					beego.Error("Fail to set message file: " + err.Error())
					return err
				}
				langData = append(langData, tempData...)
			}
			return nil
		})
		if err := i18n.SetMessageData(lang, langData); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

func ReadFile(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
