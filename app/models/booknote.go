package models

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BookNote struct {
	Id					int64			`orm:"column(id);auto"`				//The primary key id
	AuthorId			*WpUsers		`orm:"default(0);index;rel(fk);on_delete(do_nothing)"`
	ClassName			string			`orm:"size(50)"`					//enum book_note,net
	Note				string			`orm:"type(text)"`
	Source				string			`orm:"size(125)"`
	Ip					string			`orm:"size(16);type(char)"`
	Order				int				`orm:"default(1)"`
	Ding				int				`orm:"default(0)"`
	Cai					int				`orm:"default(0)"`
	Collected			int				`orm:"default(0)"`
	CreatedAt			time.Time		`orm:"auto_now_add;type(datetime)"`
	UpdatedAt			time.Time		`orm:"auto_now;type(datetime)"`
}
func (u *BookNote) TableName() string {
	return beego.AppConfig.String("MysqlPrefix") + "booknotes"
}

func init()  {
	//register models
	orm.RegisterModel(new(BookNote))
	//orm.RegisterModelWithPrefix(beego.AppConfig.String("MysqlPrefix"),new(WpUsers))
	//orm.RegisterModelWithPrefix("wp_", new(WpUsers),  new(WpUsermeta),  new(WpPosts), new(WpPostmeta), new(WpOptions))
}