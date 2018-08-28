package routers

import (
	"windigniter.com/app/controllers"
	"github.com/astaxie/beego"
	"encoding/gob"
	"time"
)

func init() {
	//register Interface
	gob.Register(controllers.SessionUser{})
	gob.Register(time.Time{})

	//load page exception
	beego.ErrorController(&controllers.ErrorsController{})

	//route map
	beego.Router("/", &controllers.MainController{}, "get:Index")
    //beego.Router("/", &controllers.JapanNewsController{}, "get:Index")
	//beego.Router("/login", &controllers.UserController{}, "get:LoginForm;post:LoginPost")
	beego.Router("/login", &controllers.UserController{}, "get,post:Login")
	beego.Router("/logout", &controllers.UserController{}, "post:Logout")
	beego.Router("/register", &controllers.UserController{}, "*:Register")
	beego.Router("/lost-password", &controllers.UserController{}, "get,post:LostPassword")
	beego.Router("/reset-password", &controllers.UserController{}, "get,post:ResetPassword")
	beego.Router("/member/u_:username([\\w]+)", &controllers.MemberController{}, "get:Index")

	beego.Router("/blog/index", &controllers.BlogController{}, "*:Index")
	beego.Router("/blog/archives/:id:int", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:year([0-9]{4})/:month([0-9]{2})/:day([0-9]{2})/:hour([0-9]{2})/:minute([0-9]{2})/:second([0-9]{2})/:postName", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:year([0-9]{4})/:month([0-9]{2})/:day([0-9]{2})/:hour([0-9]{2})/:minute([0-9]{2})/:postName", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:year([0-9]{4})/:month([0-9]{2})/:day([0-9]{2})/:hour([0-9]{2})/:postName", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:year([0-9]{4})/:month([0-9]{2})/:day([0-9]{2})/:postName", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:year([0-9]{4})/:month([0-9]{2})/:postName", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:year([0-9]{4})/:postName:string", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog/:postName:string", &controllers.BlogController{}, "get:Show")
	beego.Router("/blog", &controllers.BlogController{}, "get:Index")


	beego.Router("/japannews/content", &controllers.JapanNewsController{}, "*:NewsContent")

	beego.Router("/japan", &controllers.JapanNewsController{}, "get:Index")
	beego.Router("/japannews", &controllers.JapanNewsController{}, "get:Index")
	beego.Router("/japan/news", &controllers.JapanNewsController{}, "get:Index")
	beego.Router("/japannews/index", &controllers.JapanNewsController{}, "get:Index")
	beego.Router("/japan/news/index", &controllers.JapanNewsController{}, "get:Index")
	beego.Router("/japannews/:id:int", &controllers.JapanNewsController{}, "*:Show")
	beego.Router("/japan/news/:id:int", &controllers.JapanNewsController{}, "*:Show")
	beego.Router("/japannews/show", &controllers.JapanNewsController{}, "*:Show")
	beego.AutoRouter(&controllers.JapanNewsController{})
	beego.AutoRouter(&controllers.MemberController{})
    //beego.Include(&controllers.UserController{})
}
