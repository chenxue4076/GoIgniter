package controllers

import (
	"net/http"
	"github.com/astaxie/beego/validation"
	"log"
	"strings"
	"windigniter.com/app/services"
	"github.com/astaxie/beego"
	"fmt"
	"windigniter.com/app/libraries"
	"io/ioutil"
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
	c.LayoutSections["Scripts"] = ""
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
				//e := validation.Error{Message:Translate(lang, err.Error()), Key:key}
				valid.SetError(key, Translate(lang, err.Error()))
				c.Data["Error"] = valid.Errors
				if isAjax {
					c.Data["json"] = result
					c.ServeJSON()
				}
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
				if isAjax {
					c.Data["json"] = result
					c.ServeJSON()
				} else {
					http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, refer, 302)
				}
			}
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
	c.LayoutSections["HtmlFoot"] = ""
	lang := c.CurrentLang()//Post Data deal
	isAjax :=c.Ctx.Input.IsAjax()
	if c.Ctx.Request.Method == http.MethodPost { //POST Login deal
		var result JsonOut
		hasError := false
		username := strings.TrimSpace(c.Input().Get("username"))
		valid := validation.Validation{}
		valid.Required(username, "username").Message(Translate(lang, "user.usernameRequired"))
		if valid.HasErrors() {
			hasError = true
			var e *validation.Error
			for index, err := range valid.Errors {
				if index == 0 {
					e = err
				}
			}
			if isAjax {
				c.Data["json"] = JsonOut{true, JsonMessage{e.Message, e.Key}, ""}
				c.ServeJSON()
			}
			c.Data["Error"] = valid.Errors
		}
		//form verify success
		if ! hasError {
			//whether has this user
			db := new(services.WpUsersService)
			user, err := db.ExistUser(username)
			if err != nil {
				hasError = true
				result = JsonOut{true, JsonMessage{Translate(lang, err.Error()), "username"}, ""}
				valid.SetError("username", Translate(lang, err.Error()))
				//e := validation.Error{Message:Translate(lang, "user.userNotExist"), Key:"username"}
				c.Data["Error"] = valid.Errors
			}
			//get user info success
			if ! hasError {
				//user exists, send an email to this user
				subject := Translate(lang, "user.resetPassword") + " - " + Translate(lang, "common.siteName")
				mailBody, e := ioutil.ReadFile("resources/lang/"+lang+"/mail/resetpassword.html")
				if e != nil {
					hasError = true
					result = JsonOut{true, JsonMessage{Translate(lang, e.Error()), "username"}, ""}
					valid.SetError("", Translate(lang, e.Error()))
					c.Data["Error"] = valid.Errors
				}
				if ! hasError {
					//replace var info
					//Scan

					err := libraries.SendMail(user.UserEmail, subject, string(mailBody), "html")
					if err != nil {
						hasError = true
						result = JsonOut{true, JsonMessage{Translate(lang, e.Error()), "username"}, ""}
						valid.SetError("", Translate(lang, e.Error()))
						c.Data["Error"] = valid.Errors
					}
					if ! hasError {
						result = JsonOut{false, JsonMessage{Translate(lang, "user.emailHasSend"), "username"}, ""}
						c.Data["Success"] = Translate(lang, "user.emailHasSend")
					}
				}

			}
		}
		if isAjax {
			c.Data["json"] = result
			c.ServeJSON()
		}
	}
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