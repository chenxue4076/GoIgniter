package controllers

import (
	"strconv"

	"windigniter.com/app/libraries"
	"windigniter.com/app/services"
)

type BookNoteController struct {
	BaseController
}

// book note list
func (c *BookNoteController) Index() {
	perPage := 20
	lang := c.CurrentLang()
	pageString := c.GetString("page")
	if pageString == "" {
		pageString = "1"
	}
	page, _ := strconv.Atoi(pageString)
	db := new(services.BookNoteService)
	bookNoteList, total, err := db.BookNoteList(perPage, page)
	if err != nil {
		c.Data["Title"] = Translate(lang, "common.unknownError")
		c.Data["Content"] = Translate(lang, err.Error())
		c.Abort("Normal")
	}
	c.Data["Total"] = total
	c.Data["BookNoteList"] = bookNoteList
	c.Data["Title"] = Translate(lang, "booknote.title")
	c.Data["Pagination"] = libraries.PageList(total, page, perPage, c.Ctx.Request.URL.Path, 3, lang)
}
