package routers

import (
	"windigniter.com/app/controllers"
	"github.com/astaxie/beego"
	)

func init() {
    beego.Router("/", &controllers.MainController{})
	//beego.Router("/login", &controllers.UserController{}, "get:LoginForm;post:LoginPost")
	beego.Router("/login", &controllers.UserController{}, "get,post:Login")
	beego.Router("/register", &controllers.UserController{}, "*:Register")
	beego.Router("/member/u_:username([\\w]+)", &controllers.MemberController{}, "get:Index")
    beego.AutoRouter(&controllers.MemberController{})
    //beego.Include(&controllers.UserController{})
}
