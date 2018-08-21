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
	"strconv"
	"time"
	"golang.org/x/crypto/bcrypt"
	"errors"
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
			user, resetKey, err := db.DoResetPassword(username)
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
					link := "http://" + c.Ctx.Request.Host + "/reset-password?key=" + resetKey + "&user="+user.UserLogin
					fmt.Println("user reset password link :", link)
					//send email
					mailBodyString := fmt.Sprintf(string(mailBody), Translate(lang, "common.siteName"), link, link, Translate(lang, "common.siteName"))
					err := libraries.SendMail(user.UserEmail, subject, mailBodyString, "html")
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
	c.LayoutSections["HtmlFoot"] = ""
	lang := c.CurrentLang()
	isAjax :=c.Ctx.Input.IsAjax()
	//Post Data deal
	var username, key, password string
	valid := validation.Validation{}
	if c.Ctx.Request.Method == http.MethodPost { //POST Login deal
		username = strings.TrimSpace(c.Input().Get("username"))
		key = strings.TrimSpace(c.Input().Get("key"))
		password = strings.TrimSpace(c.Input().Get("password"))
		passwordConfirm := strings.TrimSpace(c.Input().Get("passwordConfirm"))
		valid.Required(username, "").Message(Translate(lang, "common.invalidRequest"))
		valid.Required(key, "").Message(Translate(lang, "common.invalidRequest"))
		valid.Required(password, "password").Message(Translate(lang, "user.passwordRequired"))
		valid.Required(passwordConfirm, "passwordConfirm").Message(Translate(lang, "user.passwordRequired"))
		if password != passwordConfirm {
			valid.SetError("passwordConfirm", Translate(lang, "user.passwordNotMatch"))
		}
		if valid.HasErrors() {
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
	} else {
		username = c.GetString("user")
		key = c.GetString("key")
		//invalid input
		if username == "" || key == "" {
			if isAjax {
				c.Data["json"] = JsonOut{true, JsonMessage{Translate(lang, "common.invalidRequest"), ""}, ""}
				c.ServeJSON()
			}
			c.Data["Title"] = Translate(lang,"common.invalidRequest")
			c.Data["Content"] = Translate(lang,"common.serverDealError")
			c.Abort("Normal")
		}
		c.Data["Key"] = key
		c.Data["Username"] = username
	}
	//get user info
	db := new(services.WpUsersService)
	user, err := db.ExistUser(username)
	if err != nil {
		if isAjax {
			c.Data["json"] = JsonOut{true, JsonMessage{Translate(lang, err.Error()), ""}, ""}
			c.ServeJSON()
		}
		if c.Ctx.Request.Method == http.MethodPost {
			valid.SetError("", Translate(lang, err.Error()))
			c.Data["Error"] = valid.Errors
		} else {
			c.Data["Title"] = Translate(lang,"common.invalidRequest")
			c.Data["Content"] = Translate(lang,"common.serverDealError")
			c.Abort("Normal")
		}
	} else {
		if c.Ctx.Request.Method == http.MethodPost {
			//verify key
			timeKey := strings.Split(user.UserActivationKey, ":")
			timeString, hashActivationKey := timeKey[0], timeKey[1]
			timeInt, _ := strconv.Atoi(timeString)
			now, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
			if now - timeInt > 3600 * 24 {
				if isAjax {
					c.Data["json"] = JsonOut{true, JsonMessage{Translate(lang, "user.verificationInformationHasExpired"), ""}, ""}
					c.ServeJSON()
				}
				valid.SetError("", Translate(lang, "user.verificationInformationHasExpired"))
				c.Data["Error"] = valid.Errors
			} else {
				pwderr := bcrypt.CompareHashAndPassword([]byte(hashActivationKey), []byte(key))
				if pwderr != nil {
					if libraries.HashError(pwderr).Error() == "common.hashErrMismatchedHashAndPassword" {
						pwderr = errors.New("user.resetKeyError")
					}
					if isAjax {
						c.Data["json"] = JsonOut{true, JsonMessage{Translate(lang, libraries.HashError(pwderr).Error()), ""}, ""}
						c.ServeJSON()
					}
					valid.SetError("", Translate(lang, libraries.HashError(pwderr).Error()))
					c.Data["Error"] = valid.Errors
				} else {
					//set user password
					newPassword, e :=bcrypt.GenerateFromPassword([]byte(password), 0)
					if e != nil {
						if isAjax {
							c.Data["json"] = JsonOut{true, JsonMessage{Translate(lang, libraries.HashError(e).Error()), ""}, ""}
							c.ServeJSON()
						}
						valid.SetError("", Translate(lang, libraries.HashError(e).Error()))
						c.Data["Error"] = valid.Errors
					} else {
						user.UserPass = string(newPassword)
						user.UserActivationKey = ""
						err := db.SaveUser(user, "UserPass", "UserActivationKey")
						if err != nil {
							if isAjax {
								c.Data["json"] = JsonOut{true, JsonMessage{Translate(lang, libraries.DbError(err).Error()), ""}, ""}
								c.ServeJSON()
							}
							valid.SetError("", Translate(lang, libraries.DbError(err).Error()))
							c.Data["Error"] = valid.Errors
						} else {
							refer := c.Input().Get("refer")
							if isAjax {
								c.Data["json"] = JsonOut{false, JsonMessage{Translate(lang, "common.success"), ""}, refer}
								c.ServeJSON()
							}
							c.Data["Success"] = Translate(lang, "common.success")
							http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, refer, 302)
						}
					}
				}
			}
		}
	}
	c.Data["Refer"] = beego.URLFor("UserController.Login")
	c.Data["Title"] = Translate(lang,"user.resetPassword")
	c.Data["User"] = user
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