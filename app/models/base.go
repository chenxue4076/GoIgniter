package models

import (
	"time"
)

type Users struct {
	Id					int64			`orm:"column(ID);auto"`				//The primary key id
	UserLogin			string			`orm:"size(60)"`					//The user login name
	UserPass			string			`orm:"size(255)"`
	UserNicename		string			`orm:"size(50)"`
	UserEmail			string			`orm:"size(100)"`
	UserUrl				string			`orm:"size(100)"`
	UserRegistered		time.Time		`orm:"auto_now_add;type(datetime)"`
	UserActivationKey	string
	UserStatus			int				`orm:"default(0)"`
	DisplayName			string
	Spam				int8			`orm:"default(0)"`
	Deleted				int8			`orm:"default(0)"`
}