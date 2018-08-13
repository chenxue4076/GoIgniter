package routers

import (
	"windigniter.com/app/controllers"
	"github.com/astaxie/beego"
	)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.UserController{}, "get:LoginForm;post:LoginPost")
    //beego.AutoRouter(&controllers.UserController{})
    //beego.Include(&controllers.UserController{})
}
