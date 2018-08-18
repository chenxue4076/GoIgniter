package controllers

import (
	"net/http"
	"github.com/astaxie/beego/validation"
	"log"
	"strings"
	"windigniter.com/app/services"
	"fmt"
	"github.com/astaxie/beego"
)

type UserController struct {
	BaseController
}

/*type UserLogin struct {
	Username string
	Password string
}*/

func (c *UserController) Login() {
	if c.GetSession("userInfo") != nil {
		http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, beego.URLFor("MemberController.Center"), 302)
	}
	c.LayoutSections["HtmlFoot"] = ""
	lang := c.CurrentLang()
	isAjax :=c.Ctx.Input.IsAjax()

	//Post Data deal
	if c.Ctx.Request.Method == http.MethodPost {	//POST Login deal
		refer := c.Input().Get("refer")
		username := strings.TrimSpace(c.Input().Get("username"))
		password := strings.TrimSpace(c.Input().Get("password"))
		remember := c.Input().Get("remember")
		valid := validation.Validation{}
		valid.Required(username, "username").Message(Translate(lang, "user.usernameRequired"))
		valid.Required(password, "password").Message(Translate(lang, "user.passwordRequired"))
		if valid.HasErrors() {
			var e *validation.Error
			for index, err := range valid.Errors {
				if index == 0 {
					e = err
				}
				log.Println(err.Key, err.Message)
			}
			if isAjax {
				c.Data["json"] = JsonOut{true, JsonMessage{e.Message, e.Key}, ""}
				c.ServeJSON()
			}
			c.Data["Error"] = valid.Errors
		} else {
			//TODO login action
			db := new(services.WpUsersService)
			user, key, err := db.LoginCheck(username, password)
			var result JsonOut
			if err != nil {
				fmt.Println("no info")
				result = JsonOut{true, JsonMessage{Translate(lang, err.Error()), key}, ""}
				e := validation.Error{Message:Translate(lang, err.Error()), Key:key}
				c.Data["Error"] = e
			} else {
				fmt.Println("has info")
				result = JsonOut{false, JsonMessage{Translate(lang, "user.loginSuccess"), key}, refer}
				//save session for user
				//fmt.Println(reflect.TypeOf(user.UserRegistered), reflect.ValueOf(user.UserRegistered).Kind())
				sessionUser := SessionUser{user.Id, user.UserLogin, user.UserNicename, user.UserEmail, user.UserRegistered, user.DisplayName}
				c.SetSession("userInfo", sessionUser)
				if remember != "" {
					c.Ctx.SetCookie(beego.AppConfig.String("SessionName"), c.Ctx.Input.CruSession.SessionID(), 30*24*3600)	//存储30天
				}
				/*skey := reflect.TypeOf(sessionUser)
				sValue := reflect.ValueOf(sessionUser)
				for k := 0; k < skey.NumField(); k++ {
					c.SetSession(skey.Field(k).Name, sValue.Field(k).Interface())
				}*/
			}
			if isAjax {
				c.Data["json"] = result
				c.ServeJSON()
			}
			http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, refer, 302)
		}
		//if len(username) < 6
	}
	refer := c.GetString("refer")
	c.Data["Title"] = Translate(lang,"user.login")
	c.Data["Refer"] = refer
}

func (c *UserController) Register()  {
	lang := c.CurrentLang()
	//isAjax :=c.Ctx.Input.IsAjax()
	c.Data["Title"] = Translate(lang,"user.register")
}

func (c *UserController) LostPassword() {
	lang := c.CurrentLang()
	c.Data["Title"] = Translate(lang,"user.lostPassword")
}

func (c *UserController) ResetPassword() {
	lang := c.CurrentLang()
	c.Data["Title"] = Translate(lang,"user.resetPassword")
}

func (c *UserController) Logout() {
	isAjax :=c.Ctx.Input.IsAjax()
	//Post Data deal
	if c.Ctx.Request.Method == http.MethodPost { //POST Login deal
		lang := c.CurrentLang()
		refer := c.Input().Get("refer")
		//清空session
		c.DelSession("userInfo")
		if refer == "" {
			refer = beego.URLFor("MainController.Index")
		}
		if isAjax {
			c.Data["json"] = JsonOut{false, JsonMessage{Translate(lang, "user.logoutSuccess"), ""}, refer}
			c.ServeJSON()
		}
		http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, refer, 302)
	}
	http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, beego.URLFor("MainController.Index"), 302)
}

/*func (c *UserController) LoginPost() {
	c.TplName = "user/login"
}*/