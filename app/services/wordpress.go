package services

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"windigniter.com/app/models"
	"fmt"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"windigniter.com/app/libraries"
	"time"
	"strconv"
)

type WpUsersService struct {
	BaseService
}

var o orm.Ormer

func init()  {
	o = orm.NewOrm()
	o.Using("default")
}

//wordpress options
func (s *WpUsersService) Options(optionName string) (optionValue string, err error) {
	options := models.WpOptions{OptionName:optionName}
	err = o.Read(&options, "OptionName")
	if err != nil {
		fmt.Println("has err ?", err)
		return optionValue, libraries.DbError(err)
	}
	return options.OptionValue, nil
}

// User login check
func (s *WpUsersService) LoginCheck(username, password string) (user models.WpUsers, key string, err error) {
	user, err = s.ExistUser(username)
	if err != nil {
		fmt.Println("has err ?", err)
		return user, "username", libraries.DbError(err)
	} else {
		/*gpwd, e :=bcrypt.GenerateFromPassword([]byte(password), 0)
		if e != nil {
			fmt.Println("generate password err ?", e)
		} else {
			fmt.Println("generate password ", string(gpwd))
		}*/
		//verify password
		pwderr := bcrypt.CompareHashAndPassword([]byte(user.UserPass), []byte(password))
		if pwderr != nil {
			fmt.Println("password err ?", pwderr)
			return user,"password", libraries.HashError(pwderr)
		}
		fmt.Println(user)
		return user, "username", nil
	}
	return user, "", nil
}
//wether user has exist
func (s *WpUsersService) ExistUser(username string) (user models.WpUsers, err error) {
	match, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, username)
	if match {
		user = models.WpUsers{UserEmail:username}
		err = o.Read(&user, "UserEmail")
	} else {
		user = models.WpUsers{UserLogin:username}
		err = o.Read(&user, "UserLogin")
	}
	if err != nil {
		if err == orm.ErrNoRows {
			return user, errors.New("user.userNotExist")
		} else {
			fmt.Println("has err ?", err)
			return user, libraries.DbError(err)
		}
	}
	fmt.Println("success get user ", user)
	return user, nil
}
//reset password
func (s *WpUsersService) DoResetPassword(username string) (user models.WpUsers, key string, err error) {
	user, err = s.ExistUser(username)
	if err != nil {
		return user, "", err
	}
	//更新用户重置密码字段
	key = libraries.WpGeneratePassword(20, false, false)
	//set hash for key
	hashKey, err := bcrypt.GenerateFromPassword([]byte(key), 8)
	if err != nil {
		return user, "", err
	}
	timeUnix := time.Now().Unix()
	user.UserActivationKey = strconv.FormatInt(timeUnix, 10) + ":" +  string(hashKey)
	if _, err := o.Update(&user, "UserActivationKey"); err != nil {
		return user, "", libraries.DbError(err)
	}
	return user, key,nil
}
//update user info
func (s *WpUsersService) SaveUser(user models.WpUsers, cols ...string) error {
	if _, err := o.Update(&user, cols...); err != nil {
		return err
	}
	return nil
}


// blog new list
func (s *WpUsersService) BlogList(limit, page int) (blogs []*models.WpPosts, total int64, err error) {
	wpPosts := models.WpPosts{}
	qs := o.QueryTable(wpPosts.TableName())
	total, err = qs.Count()
	if err != nil {
		fmt.Println("has err get total ?", err)
		return blogs, total, libraries.DbError(err)
	}
	offset := (page - 1) * limit
	_, err = qs.OrderBy("-ID").Limit(limit).Offset(offset).All(&blogs)
	if err != nil {
		fmt.Println("has err ?", err)
		return blogs, total, libraries.DbError(err)
	}
	return blogs, total, nil
}
// blog detail by id
func (s *WpUsersService) BlogDetail(Id int64, postName string) (blog models.WpPosts, err error) {
	//wpUser := models.WpUsers{}
	wpPosts := models.WpPosts{}
	qs := o.QueryTable(wpPosts.TableName())
	fmt.Println("BlogDetail", Id, postName)
	if Id != 0 {
		qs = qs.Filter("ID", Id)
	} else if postName != "" {
		qs = qs.Filter("post_name", postName)
	} else {
		return blog, errors.New("common.ErrArgs")
	}
	err = qs.Filter("post_status", "publish").RelatedSel().One(&blog)
	if err != nil {
		return blog, err
	}
	return blog, nil
}