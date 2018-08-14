package main

import (
	_ "windigniter.com/routers"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
		"html"
)

func main() {
	//fmt.Println(os.Getenv("BEEGO_RUNMODE"))
	//StaticDir["/public"] = "public"
	beego.SetStaticPath("/public", "public")
	beego.SetStaticPath("/favicon.ico", "public/favicon.ico")

	beego.AddFuncMap("i18n", i18n.Tr)
	beego.AddFuncMap("html", html.UnescapeString)

	beego.Run()
}

