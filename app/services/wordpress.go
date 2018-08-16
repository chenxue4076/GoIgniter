package services

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"windigniter.com/app/models"
	"fmt"
	"regexp"
)

type WpUsersService struct {
	BaseService
}

var o orm.Ormer

func init()  {
	o = orm.NewOrm()
	o.Using("default")
}

func (s *WpUsersService) LoginCheck(username, password string) error {
	//\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*
	match, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, username)
	var user models.WpUsers
	if match {
		fmt.Println("col email")
		user = models.WpUsers{UserEmail:username}
	} else {
		fmt.Println("col name")
		user = models.WpUsers{UserLogin:username}
	}

	err := o.Read(&user)
	/*if err == orm.ErrNoRows {
		fmt.Println("not found", err)
		return false
	} else if err == orm.ErrMissPK {
		fmt.Println("not found key", err)
		return false
	} else {
		fmt.Println("has err ?", err)
		fmt.Println(err)
		return true
	}*/
	if err != nil {
		fmt.Println("has err ?", err)
		return err
	} else {
		fmt.Println(user)
		return nil
	}
	return nil
}
