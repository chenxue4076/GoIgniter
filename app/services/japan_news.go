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
// japan news list
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
//japan news detail
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
//japan easy news detail
func (s *JapanNewsService) JapanEasyNewsDetail(newsId string) (news models.JapanEasyNews, err error) {
	easyNews := models.JapanEasyNews{}
	err = o.QueryTable(easyNews.TableName()).Filter("NewsId", newsId).One(&news)
	if err != nil {
		return news, libraries.DbError(err)
	}
	return news, nil
}
// japan easy news insert
func (s *JapanNewsService) SaveEasyNews(news models.JapanEasyNews) (id int, err error)  {
	id64, err := o.Insert(&news)
	if err != nil {
		return 0, libraries.DbError(err)
	}
	id, _ = strconv.Atoi(strconv.FormatInt(id64, 10))
	return id, err
}
// japan easy news update
func (s *JapanNewsService) UpdateEasyNews(news models.JapanEasyNews, cols ...string) error {
	if _, err := o.Update(&news, cols...); err != nil {
		return err
	}
	return nil
}
