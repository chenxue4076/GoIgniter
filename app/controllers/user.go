package controllers

import (
	"net/http"
	"github.com/astaxie/beego/validation"
	"log"
	"strings"
	"windigniter.com/app/services"
	"fmt"
)

type UserController struct {
	BaseController
}

/*type UserLogin struct {
	Username string
	Password string
}*/

func (c *UserController) Login() {
	c.LayoutSections["HtmlFoot"] = ""
	lang := c.CurrentLang()
	isAjax :=c.Ctx.Input.IsAjax()

	//Post Data deal
	if c.Ctx.Request.Method == http.MethodPost {	//POST Login deal
		refer := c.Input().Get("refer")
		username := strings.TrimSpace(c.Input().Get("username"))
		password := strings.TrimSpace(c.Input().Get("password"))
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
			hasUser := db.LoginCheck(username, password)
			if hasUser {
				fmt.Println("查询到了")
			} else {
				fmt.Println("查询不到呢")
			}
			if isAjax {
				c.Data["json"] = JsonOut{false, JsonMessage{Translate(lang, "user.loginSuccess"), ""}, refer}
				c.ServeJSON()
			}
			//http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, refer, 302)
		}
		//if len(username) < 6
	}
	refer := c.GetString("refer")
	c.Data["Title"] = Translate(lang,"user.login")
	c.Data["Refer"] = refer
	//c.LayoutSections["HeaderMeta"] = "user/headermeta.html"
	//c.LayoutSections["HtmlHead"] = ""
	//c.LayoutSections["Scripts"] = "user/scripts.html"
}

/*func (c *UserController) LoginPost() {
	c.TplName = "user/login"
}*/