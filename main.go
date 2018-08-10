package main

import (
	_ "windigniter.com/routers"
	"github.com/astaxie/beego"
)

func main() {
	//StaticDir["/public"] = "public"
	beego.SetStaticPath("/public", "public")
	beego.SetStaticPath("/favicon.ico", "public/favicon.ico")
	beego.Run()
}

