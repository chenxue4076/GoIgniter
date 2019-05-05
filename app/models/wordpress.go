package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"time"
)

var preTable = "wp_"

type WpUsers struct {
	Id                int64     `orm:"column(ID);auto"` //The primary key id
	UserLogin         string    `orm:"size(60);unique"` //The user login name
	UserPass          string    `orm:"size(255)"`
	UserNicename      string    `orm:"size(50);index"`
	UserEmail         string    `orm:"size(100);unique"`
	UserUrl           string    `orm:"size(100)"`
	UserRegistered    time.Time `orm:"auto_now_add;type(datetime)"`
	UserActivationKey string
	UserStatus        int `orm:"default(0)"`
	DisplayName       string
	//Spam				int8			`orm:"default(0)"`
	//Deleted				int8			`orm:"default(0)"`
	UserPosts    []*WpPosts    `orm:"reverse(many)"`
	UserMetas    []*WpUsermeta `orm:"reverse(many)"`
	UserBookNote []*BookNote   `orm:"reverse(many)"`
}

func (u *WpUsers) TableName() string {
	return preTable + "users"
}

type WpUsermeta struct {
	UmetaId   int64    `orm:"auto"`
	UserId    *WpUsers `orm:"default(0);index;rel(fk);on_delete(cascade)"`
	MetaKey   string   `orm:"size(255);null;index"`
	MetaValue string   `orm:"type(text);null"`
}

func (u *WpUsermeta) TableName() string {
	return preTable + "usermeta"
}

type WpPosts struct {
	Id                  int64                  `orm:"column(ID);auto"`
	PostAuthor          *WpUsers               `orm:"column(post_author);default(0);index;rel(fk);on_delete(do_nothing)"`
	PostDate            time.Time              `orm:"auto_now_add;type(datetime)"`
	PostDateGmt         time.Time              `orm:"auto_now_add;type(datetime)"`
	PostContent         string                 `orm:"type(text)"`
	PostTitle           string                 `orm:"type(text)"`
	PostExcerpt         string                 `orm:"type(text)"`
	PostStatus          string                 `orm:"size(20);default(publish)"`
	CommentStatus       string                 `orm:"size(20);default(open)"`
	PingStatus          string                 `orm:"size(20);default(open)"`
	PostPassword        string                 `orm:"size(255)"`
	PostName            string                 `orm:"size(200)"`
	ToPing              string                 `orm:"type(text)"`
	Pinged              string                 `orm:"type(text)"`
	PostModified        time.Time              `orm:"auto_now;type(datetime)"`
	PostModifiedGmt     time.Time              `orm:"auto_now;type(datetime)"`
	PostContentFiltered string                 `orm:"type(text)"`
	PostParent          int64                  `orm:"default(0);index"`
	Guid                string                 `orm:"size(255)"`
	MenuOrder           int                    `orm:"default(0)"`
	PostType            string                 `orm:"size(20)"`
	PostMimeType        string                 `orm:"size(100)"`
	CommentCount        int64                  `orm:"default(0)"`
	PostMeta            []*WpPostmeta          `orm:"reverse(many)"`
	PostTags            []*WpTermRelationships `orm:"reverse(many)"`
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
	MetaId    int64    `orm:"auto"`
	PostId    *WpPosts `orm:"default(0);index;rel(fk);on_delete(cascade)"`
	MetaKey   string   `orm:"size(255);null;index"`
	MetaValue string   `orm:"type(text)"`
}

func (u *WpPostmeta) TableName() string {
	return preTable + "postmeta"
}

type WpTerms struct {
	TermId    int64         `orm:"auto"`
	Name      string        `orm:"size(200);index"`
	Slug      string        `orm:"size(200);index"`
	TermGroup int64         `orm:"default(0)"`
	TermMate  []*WpTermmeta `orm:"reverse(many)"`
}

func (u *WpTerms) TableName() string {
	return preTable + "terms"
}

type WpTermmeta struct {
	MetaId    int64    `orm:"auto"`
	TermId    *WpTerms `orm:"default(0);index;rel(fk);on_delete(cascade)"`
	MetaKey   string   `orm:"size(255);null;index"`
	MetaValue string   `orm:"type(text)"`
}

func (u *WpTermmeta) TableName() string {
	return preTable + "termmeta"
}

type WpTermTaxonomy struct {
	TermTaxonomyId int64  `orm:"auto"`
	TermId         int64  `orm:"default(0)"`
	Taxonomy       string `orm:"size(32);index"`
	Description    string `orm:"type(text)"`
	Parent         int64  `orm:"default(0)"`
	Count          int64  `orm:"default(0)"`
}

func (u *WpTermTaxonomy) TableName() string {
	return preTable + "term_taxonomy"
}
func (u *WpTermTaxonomy) TableIndex() [][]string {
	return [][]string{
		[]string{"TermId", "Taxonomy"},
	}
}

type WpTermRelationships struct {
	Id             int64    `orm:"auto"`
	ObjectId       *WpPosts `orm:"default(0);rel(fk);on_delete(cascade)"`
	TermTaxonomyId int64    `orm:"default(0);index"`
	TermOrder      int64    `orm:"default(0)"`
}

func (u *WpTermRelationships) TableName() string {
	return preTable + "term_relationships"
}
func (u *WpTermRelationships) TableIndex() [][]string {
	return [][]string{
		[]string{"ObjectId", "TermTaxonomyId"},
	}
}

type WpOptions struct {
	OptionId    int64  `orm:"auto"`
	OptionName  string `orm:"size(191);index"`
	OptionValue string `orm:"type(text)"`
	Autoload    string `orm:"size(20);default(yes)"`
}

func (u *WpOptions) TableName() string {
	return preTable + "options"
}

//不是创建表，只是用来接收查询结果，
type Tags struct {
	TermId int64
	Name   string
	Slug   string
	Count  int64
}

func init() {
	//register models
	orm.RegisterModel(new(WpUsers), new(WpUsermeta), new(WpPosts), new(WpPostmeta), new(WpTerms), new(WpTermmeta), new(WpTermTaxonomy), new(WpTermRelationships), new(WpOptions))
	//orm.RegisterModelWithPrefix(beego.AppConfig.String("MysqlPrefix"),new(WpUsers))
	//orm.RegisterModelWithPrefix("wp_", new(WpUsers),  new(WpUsermeta),  new(WpPosts), new(WpPostmeta), new(WpOptions))
}
