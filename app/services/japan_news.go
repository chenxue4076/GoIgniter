package services

import (
	"windigniter.com/app/models"
	"fmt"
	"windigniter.com/app/libraries"
	"strconv"
)

type JapanNewsService struct {
	BaseService
}

/*func init()  {
	o = orm.NewOrm()
	o.Using("default")
}*/

func (s *JapanNewsService) JapanNewsList(limit, page, status int) (newsList []*models.JapanNews, total int, err error) {
	japanNews := models.JapanNews{}
	qs := o.QueryTable(japanNews.TableName())
	if status > -1 {
		qs = qs.Filter("status", 1)
	}
	total64, err := qs.Count()
	total, _ = strconv.Atoi(strconv.FormatInt(total64, 10))
	if err != nil {
		fmt.Println("has err get total ?", err)
		return newsList, total, libraries.DbError(err)
	}
	offset := (page - 1) * limit
	_, err = qs.OrderBy("-ID").Limit(limit).Offset(offset).All(&newsList)
	if err != nil {
		fmt.Println("has err ?", err)
		return newsList, total, libraries.DbError(err)
	}
	return newsList, total, nil
}

func (s *JapanNewsService) JapanNewsDetail(Id int64) (news models.JapanNews, err error) {
	//wpUser := models.WpUsers{}
	japanNews := models.JapanNews{}
	qs := o.QueryTable(japanNews.TableName()).Filter("ID", Id)
	err = qs.One(&news)
	if err != nil {
		return news, libraries.DbError(err)
	}
	return news, nil
}