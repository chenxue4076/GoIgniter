package controllers

import (
	"github.com/astaxie/beego"
		"os"
		"strings"
	"io/ioutil"
	"html/template"
	"path/filepath"
	"github.com/beego/i18n"
	"sync"
)

var once sync.Once

type BaseController struct {
	beego.Controller
}

// all controllers init
func (c *BaseController) Prepare()  {
	//page params
	c.TplExt = "html"
	controller, action := c.GetControllerAndAction()
	c.TplName = strings.Replace(strings.ToLower(controller), "controller", "", -1) + "/" + strings.ToLower(action) + "." +  c.TplExt
	c.Layout = "layout/common."+c.TplExt
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HeaderMeta"] = ""
	c.LayoutSections["HtmlHead"] = ""
	c.LayoutSections["HtmlFoot"] = "layout/footer."+c.TplExt
	c.LayoutSections["Scripts"] = ""
	c.LayoutSections["SideBar"] = ""


	c.Data["XsrfData"] = template.HTML(c.XSRFFormHTML())
	c.Data["Lang"] = "zh-CN"
	once.Do(loadLangs)

}

func loadLangs()  {
	//language choose
	langs := []string {"zh-CN", "en-US"}
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
