package models

import (
	"time"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/astaxie/beego/orm"
)

type WpUsers struct {
	Id					int64			`orm:"column(ID);auto"`				//The primary key id
	UserLogin			string			`orm:"size(60);unique"`				//The user login name
	UserPass			string			`orm:"size(255)"`
	UserNicename		string			`orm:"size(50)"`
	UserEmail			string			`orm:"size(100);unique"`
	UserUrl				string			`orm:"size(100)"`
	UserRegistered		time.Time		`orm:"auto_now_add;type(datetime)"`
	UserActivationKey	string
	UserStatus			int				`orm:"default(0)"`
	DisplayName			string
	//Spam				int8			`orm:"default(0)"`
	//Deleted				int8			`orm:"default(0)"`
}

func init()  {
	//register models
	orm.RegisterModel(new(WpUsers))
	//orm.RegisterModelWithPrefix(beego.AppConfig.String("MysqlPrefix"),new(WpUsers))
	//orm.RegisterModelWithPrefix("wp_",new(WpUsers))
}