package main

import (
	_ "windigniter.com/routers"
	"github.com/astaxie/beego"
	"os"
	"fmt"
)

func main() {
	/*err := os.Setenv("BEEGO_RUNMODE", "dev")
	if err != nil {
		fmt.Println(err.Error())
		return
	}*/
	fmt.Println(os.Getenv("BEEGO_RUNMODE"))
	//StaticDir["/public"] = "public"
	beego.SetStaticPath("/public", "public")
	beego.SetStaticPath("/favicon.ico", "public/favicon.ico")
	beego.Run()
}

