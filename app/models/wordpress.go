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
	UserPosts			[]*WpPosts		`orm:"reverse(many)"`
}

type WpPosts struct {
	Id					int64			`orm:"column(ID);auto"`
	PostAuthor			*WpUsers		`orm:"default(0);index;rel(fk);on_delete(do_nothing)"`
	PostDate			time.Time		`orm:"auto_now_add;type(datetime)"`
	PostDateGmt			time.Time		`orm:"auto_now_add;type(datetime)"`
	PostContent			string			`orm:"type(text)"`
	PostTitle			string			`orm:"type(text)"`
	PostExcerpt			string			`orm:"type(text)"`
	PostStatus			string			`orm:"size(20);default(publish)"`
	CommentStatus		string			`orm:"size(20);default(open)"`
	PingStatus			string			`orm:"size(20);default(open)"`
	PostPassword		string			`orm:"size(255)"`
	PostName			string			`orm:"size(200)"`
	ToPing				string			`orm:"type(text)"`
	Pinged				string			`orm:"type(text)"`
	PostModified		time.Time		`orm:"auto_now;type(datetime)"`
	PostModifiedGmt		time.Time		`orm:"auto_now;type(datetime)"`
	PostContentFiltered	string			`orm:"type(text)"`
	PostParent			int64			`orm:"default(0);index"`
	Guid				string			`orm:"size(255)"`
	MenuOrder			int				`orm:"default(0)"`
	PostType			string			`orm:"size(20)"`
	PostMimeType		string			`orm:"size(100)"`
	CommentCount		int64			`orm:"default(0)"`
	//Author				*WpUsers		`orm:"rel(fk)"`
}


func init()  {
	//register models
	orm.RegisterModel(new(WpUsers), new(WpPosts))
	//orm.RegisterModelWithPrefix(beego.AppConfig.String("MysqlPrefix"),new(WpUsers))
	//orm.RegisterModelWithPrefix("wp_",new(WpUsers))
}