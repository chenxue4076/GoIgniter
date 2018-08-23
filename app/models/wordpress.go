package models

import (
	"time"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/astaxie/beego/orm"
)

var preTable = "wp_"

type WpUsers struct {
	Id					int64			`orm:"column(ID);auto"`				//The primary key id
	UserLogin			string			`orm:"size(60);unique"`				//The user login name
	UserPass			string			`orm:"size(255)"`
	UserNicename		string			`orm:"size(50);index"`
	UserEmail			string			`orm:"size(100);unique"`
	UserUrl				string			`orm:"size(100)"`
	UserRegistered		time.Time		`orm:"auto_now_add;type(datetime)"`
	UserActivationKey	string
	UserStatus			int				`orm:"default(0)"`
	DisplayName			string
	//Spam				int8			`orm:"default(0)"`
	//Deleted				int8			`orm:"default(0)"`
	UserPosts			[]*WpPosts		`orm:"reverse(many)"`
	UserMetas			[]*WpUsermeta	`orm:"reverse(many)"`
}
func (u *WpUsers) TableName() string {
	return preTable + "users"
}

type WpUsermeta struct {
	UmetaId					int64			`orm:"auto"`
	UserId					*WpUsers		`orm:"default(0);index;rel(fk);on_delete(cascade)"`
	MetaKey					string			`orm:"size(255);null;index"`
	MetaValue				string			`orm:"type(text);null"`
}
func (u *WpUsermeta) TableName() string {
	return preTable + "usermeta"
}

type WpPosts struct {
	Id					int64			`orm:"column(ID);auto"`
	PostAuthor			*WpUsers		`orm:"column(post_author);default(0);index;rel(fk);on_delete(do_nothing)"`
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
	PostMeta			[]*WpPostmeta	`orm:"reverse(many)"`
}
func (u *WpPosts) TableName() string {
	return preTable + "posts"
}
func (u *WpPosts) TableIndex() [][]string {
	return [][]string{
		[]string{"PostType", "PostStatus", "PostDate", "Id"},
	}
}

type WpPostmeta struct {
	MetaId					int64			`orm:"auto"`
	PostId					*WpPosts		`orm:"default(0);index;rel(fk);on_delete(cascade)"`
	MetaKey					string			`orm:"size(255);null;index"`
	MetaValue				string			`orm:"type(text)"`
}
func (u *WpPostmeta) TableName() string {
	return preTable + "postmeta"
}

type WpOptions struct {
	OptionId				int64			`orm:"auto"`
	OptionName				string			`orm:"size(191);index"`
	OptionValue				string			`orm:"type(text)"`
	Autoload				string			`orm:"size(20);default(yes)"`
}
func (u *WpOptions) TableName() string {
	return preTable + "options"
}

func init()  {
	//register models
	orm.RegisterModel(new(WpUsers),  new(WpUsermeta),  new(WpPosts), new(WpPostmeta), new(WpOptions))
	//orm.RegisterModelWithPrefix(beego.AppConfig.String("MysqlPrefix"),new(WpUsers))
	//orm.RegisterModelWithPrefix("wp_", new(WpUsers),  new(WpUsermeta),  new(WpPosts), new(WpPostmeta), new(WpOptions))
}