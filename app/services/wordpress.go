package services

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"windigniter.com/app/models"
	"fmt"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type WpUsersService struct {
	BaseService
}

var o orm.Ormer

func init()  {
	o = orm.NewOrm()
	o.Using("default")
}

func (s *WpUsersService) LoginCheck(username, password string) (user models.WpUsers, key string, err error) {
	//\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*
	match, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, username)
	//var user models.WpUsers
	//var err error
	if match {
		user = models.WpUsers{UserEmail:username}
		err = o.Read(&user, "UserEmail")
	} else {
		user = models.WpUsers{UserLogin:username}
		err = o.Read(&user, "UserLogin")
	}
	if err != nil {
		fmt.Println("has err ?", err)
		return user, "username", DbError(err)
	} else {
		/*gpwd, e :=bcrypt.GenerateFromPassword([]byte(password), 0)
		if e != nil {
			fmt.Println("generate password err ?", e)
		} else {
			fmt.Println("generate password ", string(gpwd))
		}*/
		//verify password
		pwderr := bcrypt.CompareHashAndPassword([]byte(user.UserPass), []byte(password))
		if pwderr != nil {
			fmt.Println("password err ?", pwderr)
			return user,"password", HashError(pwderr)
		}
		fmt.Println(user)
		return user, "username", nil
	}
	return user, "", nil
}

func (s *WpUsersService) ExistUser(username string) (user models.WpUsers, err error) {
	match, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, username)
	if match {
		user = models.WpUsers{UserEmail:username}
		err = o.Read(&user, "UserEmail")
	} else {
		user = models.WpUsers{UserLogin:username}
		err = o.Read(&user, "UserLogin")
	}
	if err != nil {
		if err == orm.ErrNoRows {
			return user, errors.New("user.userNotExist")
		} else {
			fmt.Println("has err ?", err)
			return user, DbError(err)
		}
	}
	fmt.Println("success get user ", user)
	return user, nil
}