package services

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/astaxie/beego"
	"errors"
	"golang.org/x/crypto/bcrypt"
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

func DbError(err error) error {
	var result string
	switch err {
		case orm.ErrNoRows:
			result = "common.ormErrNoRows"
		case orm.ErrMissPK:
			result = "common.ormErrMissPK"
		case orm.ErrTxHasBegan:
			result = "common.ormErrTxHasBegan"
		case orm.ErrTxDone:
			result = "common.ormErrTxDone"
		case orm.ErrMultiRows:
			result = "common.ormErrMultiRows"
		case orm.ErrStmtClosed:
			result = "common.ormErrStmtClosed"
		case orm.ErrArgs:
			result = "common.ormErrArgs"
		case orm.ErrNotImplement:
			result = "common.ormErrNotImplement"
		default:
			result = err.Error()
	}
	return errors.New(result)
}

func HashError(err error) error {
	var result string
	switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			result = "common.hashErrMismatchedHashAndPassword"
		case bcrypt.ErrHashTooShort:
			result = "common.hashErrHashTooShort"
		default:
			result = err.Error()
	}
	return errors.New(result)
}