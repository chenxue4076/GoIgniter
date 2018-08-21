package services

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/astaxie/beego"
)

type BaseService struct {

}

func init()  {
	//fmt.Println("service base init")
	if beego.AppConfig.String("RunMode") == beego.DEV {
		orm.Debug = true
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("MysqlUser")+":"+beego.AppConfig.String("MysqlPass")+"@tcp("+beego.AppConfig.String("MysqlHost")+":"+beego.AppConfig.String("MysqlPort")+")/"+beego.AppConfig.String("MysqlDb")+"?charset="+beego.AppConfig.String("MysqlCharSet"))
	orm.SetMaxIdleConns("default", 15)
	orm.SetMaxOpenConns("default", 30)
	//orm.DefaultTimeLoc
}

