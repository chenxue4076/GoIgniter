package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type JapanEasyNews struct {
	Id                  int64     `orm:"column(id);auto"` //The primary key id
	NewsId              string    `orm:"size(50)"`
	NewsPrearrangedTime time.Time `orm:"auto_now_add;type(datetime)"`
	Title               string    `orm:"size(255)"`
	TitleWithRuby       string    `orm:"type(text)"`
	OutlineWithRuby     string    `orm:"type(text)"`
	NewsWebImageUri     string    `orm:"size(255)"`
	NewsWebMovieUri     string    `orm:"size(255)"`
	NewsEasyVoiceUri    string    `orm:"size(255)"`
	Status              int       `orm:"default(0)"`
	CreatedAt           time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt           time.Time `orm:"auto_now;type(datetime)"`
}

func (u *JapanEasyNews) TableName() string {
	return beego.AppConfig.String("MysqlPrefix") + "japan_easy_news"
}

type JapanNews struct {
	Id           int64     `orm:"column(id);auto"` //The primary key id
	AuthorId     int64     `orm:"default(0)"`
	NewsId       string    `orm:"size(50);index"`
	Title        string    `orm:"size(255)"`
	TitleRuby    string    `orm:"type(text)"`
	DescribeRuby string    `orm:"type(text)"`
	Views        int       `orm:"default(0)"`
	Featured     string    `orm:"size(255)"`
	Media        string    `orm:"size(255)"`
	Content      string    `orm:"type(text)"`
	Dict         string    `orm:"type(text)"`
	Ding         int       `orm:"default(0)"`
	Cai          int       `orm:"default(0)"`
	Pubdate      time.Time `orm:"type(datetime)"`
	Status       int       `orm:"default(0)"`
	CreatedAt    time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt    time.Time `orm:"auto_now;type(datetime)"`
}

func (u *JapanNews) TableName() string {
	return beego.AppConfig.String("MysqlPrefix") + "japan_news"
}

func init() {
	//register models
	orm.RegisterModel(new(JapanEasyNews), new(JapanNews))
}
