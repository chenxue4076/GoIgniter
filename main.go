package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"html"
	"strings"
	"windigniter.com/app/libraries"
	_ "windigniter.com/routers"
)

func main() {
	//fmt.Println(os.Getenv("BEEGO_RUNMODE"))
	//StaticDir["/public"] = "public"
	beego.SetStaticPath("/public", "public")
	beego.SetStaticPath("/favicon.ico", "public/favicon.ico")

	beego.AddFuncMap("i18n", i18n.Tr)
	beego.AddFuncMap("html", html.UnescapeString)
	//beego.AddFuncMap("dateFormat", libraries.DateFormat)
	beego.AddFuncMap("wpUrlFormat", libraries.WordPressUrlFormat)
	//beego.AddFuncMap("removeHtml", libraries.RemoveHtml)
	beego.AddFuncMap("contains", strings.Contains)
	beego.AddFuncMap("lower", strings.ToLower)

	beego.Run()
}
