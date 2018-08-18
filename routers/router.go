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

	//route map
    beego.Router("/", &controllers.MainController{}, "get:Index")
	//beego.Router("/login", &controllers.UserController{}, "get:LoginForm;post:LoginPost")
	beego.Router("/login", &controllers.UserController{}, "get,post:Login")
	beego.Router("/logout", &controllers.UserController{}, "post:Logout")
	beego.Router("/register", &controllers.UserController{}, "*:Register")
	beego.Router("/lost-password", &controllers.UserController{}, "get,post:LostPassword")
	beego.Router("/reset-password", &controllers.UserController{}, "get,post:ResetPassword")
	beego.Router("/member/u_:username([\\w]+)", &controllers.MemberController{}, "get:Index")
    beego.AutoRouter(&controllers.MemberController{})
    //beego.Include(&controllers.UserController{})
}
