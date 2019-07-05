package services

import (
	"fmt"
	"strconv"

	"windigniter.com/app/libraries"
	"windigniter.com/app/models"
)

type BookNoteService struct {
	BaseService
}

func (s *BookNoteService) BookNoteList(limit, page int) (bookNoteList []*models.BookNote, total int, err error) {
	bookNote := models.BookNote{}
	qs := o.QueryTable(bookNote.TableName())
	total64, err := qs.Count()
	total, _ = strconv.Atoi(strconv.FormatInt(total64, 10))
	if err != nil {
		fmt.Println("has err get total ?", err)
		return bookNoteList, total, libraries.DbError(err)
	}
	offset := (page - 1) * limit
	_, err = qs.OrderBy("-ID").Limit(limit).Offset(offset).All(&bookNoteList)
	if err != nil {
		fmt.Println("has err ?", err)
		return bookNoteList, total, libraries.DbError(err)
	}
	return bookNoteList, total, nil
}
