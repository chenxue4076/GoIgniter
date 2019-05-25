package services

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"sort"
	"strconv"
	"time"
	"windigniter.com/app/libraries"
	"windigniter.com/app/models"
)

type WpUsersService struct {
	BaseService
}

//var o orm.Ormer

/*func init()  {
	o = orm.NewOrm()
	o.Using("default")
}*/

//wordpress options
func (s *WpUsersService) Options(optionName string) (optionValue string, err error) {
	options := models.WpOptions{OptionName: optionName}
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
			return user, "password", libraries.HashError(pwderr)
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
		user = models.WpUsers{UserEmail: username}
		err = o.Read(&user, "UserEmail")
	} else {
		user = models.WpUsers{UserLogin: username}
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
	user.UserActivationKey = strconv.FormatInt(timeUnix, 10) + ":" + string(hashKey)
	if _, err := o.Update(&user, "UserActivationKey"); err != nil {
		return user, "", libraries.DbError(err)
	}
	return user, key, nil
}

//update user info
func (s *WpUsersService) SaveUser(user models.WpUsers, cols ...string) error {
	if _, err := o.Update(&user, cols...); err != nil {
		return err
	}
	return nil
}

// blog new list
func (s *WpUsersService) BlogList(limit, page int) (blogs []*models.WpPosts, total int, err error) {
	wpPosts := models.WpPosts{}
	qs := o.QueryTable(wpPosts.TableName()).Filter("post_status", "publish")
	total64, err := qs.Count()
	total, _ = strconv.Atoi(strconv.FormatInt(total64, 10))
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
	qs := o.QueryTable(wpPosts.TableName()).Filter("post_status", "publish")
	//fmt.Println("BlogDetail", Id, postName)
	if Id != 0 {
		qs = qs.Filter("ID", Id)
	} else if postName != "" {
		qs = qs.Filter("post_name", postName)
	} else {
		return blog, errors.New("common.ErrArgs")
	}
	err = qs.RelatedSel().One(&blog)
	if err != nil {
		return blog, libraries.DbError(err)
	}
	return blog, nil
}

// blog show tag list
func (s *WpUsersService) Tags(Id int64, tType string) (tags []*models.WpTerms, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, err
	}
	if tType == "" {
		tType = "post_tag"
	}
	wpTerms := models.WpTerms{}
	wpTermRelationships := models.WpTermRelationships{}
	wpTermTaxonomy := models.WpTermTaxonomy{}
	sql := qb.Select("t.term_id", "t.name", "t.slug").
		From(wpTermRelationships.TableName() + " AS r").
		InnerJoin(wpTermTaxonomy.TableName() + " AS tt").On("r.term_taxonomy_id = tt.term_taxonomy_id").
		InnerJoin(wpTerms.TableName() + " AS t").On("tt.term_id = t.term_id").
		Where("r.object_id = ?").And("tt.taxonomy = ?").String()
	num, err := o.Raw(sql, Id, tType).QueryRows(&tags)
	if err != nil {
		return nil, libraries.DbError(err)
	}
	if num == 0 {
		return nil, nil
	}
	return tags, nil
}

//all tag list
func (s *WpUsersService) TagList(tType string) (tags []*models.Tags, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, err
	}
	if tType == "" {
		tType = "post_tag"
	}
	wpTerms := models.WpTerms{}
	wpTermRelationships := models.WpTermRelationships{}
	wpTermTaxonomy := models.WpTermTaxonomy{}
	sql := qb.Select("distinct t.term_id", "t.name", "t.slug", "tt.count").
		From(wpTermRelationships.TableName() + " AS r").
		InnerJoin(wpTermTaxonomy.TableName() + " AS tt").On("r.term_taxonomy_id = tt.term_taxonomy_id").
		InnerJoin(wpTerms.TableName() + " AS t").On("tt.term_id = t.term_id").
		Where("tt.taxonomy = ?").String()
	num, err := o.Raw(sql, tType).QueryRows(&tags)
	if err != nil {
		return nil, libraries.DbError(err)
	}
	if num == 0 {
		return nil, nil
	}
	return tags, nil
}

//获取terms
func (s *WpUsersService) TagInfo(slug string) (tag models.WpTerms, err error) {
	wpTerms := models.WpTerms{}
	qs := o.QueryTable(wpTerms.TableName()).Filter("slug", slug)
	err = qs.RelatedSel().One(&tag)
	if err != nil {
		return tag, libraries.DbError(err)
	}
	return tag, nil
}

// get blog by tag
func (s *WpUsersService) BlogListByTermId(termId int64, limit, page int) (blogs []*models.WpPosts, total int, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, 0, err
	}
	wpTermRelationships := models.WpTermRelationships{}
	wpTermTaxonomy := models.WpTermTaxonomy{}
	wpPosts := models.WpPosts{}
	sql := qb.Select("p.*").
		From(wpPosts.TableName() + " AS p").
		InnerJoin(wpTermRelationships.TableName() + " AS r").On("r.object_id = p.ID").
		InnerJoin(wpTermTaxonomy.TableName() + " AS tt").On("r.term_taxonomy_id = tt.term_taxonomy_id").
		Where("p.post_status = ?").And("tt.term_id = ?").String()
	num, err := o.Raw(sql, "publish", termId).QueryRows(&blogs)
	total, _ = strconv.Atoi(strconv.FormatInt(num, 10))
	if err != nil {
		return nil, 0, libraries.DbError(err)
	}
	if total == 0 {
		return nil, 0, nil
	}
	return blogs, total, nil
}

// get blog by date
func (s *WpUsersService) BlogListByDate(year string, month string, perPage int, page int) (blogs []*models.WpPosts, total int, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, 0, err
	}
	wpPosts := models.WpPosts{}
	sql := qb.Select("*").From(wpPosts.TableName()).Where("post_status = ?").And("LEFT(post_date, 7) = ?").String()
	num, err := o.Raw(sql, "publish", year+"-"+month ).QueryRows(&blogs)
	total, _ = strconv.Atoi(strconv.FormatInt(num, 10))
	if err != nil {
		return nil, 0, libraries.DbError(err)
	}
	if total == 0 {
		return nil, 0, nil
	}
	return blogs, total, nil
}

// blog ArchiveList
func (s *WpUsersService) ArchiveList() (archives []models.Archives, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, err
	}
	posts := []*models.WpPosts{}
	wpPosts := models.WpPosts{}
	sql := qb.Select("post_date").From(wpPosts.TableName()).Where("post_status = ?").String()
	num, err := o.Raw(sql, "publish").QueryRows(&posts)
	if err != nil {
		return nil, libraries.DbError(err)
	}
	if num == 0 {
		return nil, nil
	}
	mapDate := make(map[int] models.Archives)
	for _, post := range posts {
		monthItem, _ := strconv.Atoi(strconv.Itoa(post.PostDate.Year()) + fmt.Sprintf("%02d",int(post.PostDate.Month())))
		if mapDate[monthItem] != (models.Archives{}) {
			temp := mapDate[monthItem]
			temp.Count++
			mapDate[monthItem] = temp
		} else {
			mapDate[monthItem] = models.Archives{Year: post.PostDate.Year(), MonthName:post.PostDate.Month(), Month: fmt.Sprintf("%02d",int(post.PostDate.Month())), Count:1}
		}
	}
	var keys []int
	for k := range mapDate {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	//fmt.Println(keys)
	for _, k := range keys {
		archives = append(archives, mapDate[k])
	}
	return archives, nil


}