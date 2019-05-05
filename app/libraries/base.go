package libraries

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//var langs = []string {"zh-CN", "en-US"}
// get dir names at this path
func CurrentDirs(path string) (dirs []string, err error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, fi.Name())
		}
	}
	//fmt.Println("lang dir : ", dirs)
	//fmt.Println("lang static dir : ", langs)
	return dirs, nil
}

// load language info files
func LoadLangs() {
	langs, _ := CurrentDirs("resources/lang")
	if len(langs) < 1 {
		fmt.Println("fail to get language package")
		return
	}
	//language choose
	//langs := []string {"zh-CN", "en-US"}
	for _, lang := range langs {
		langData := make([]byte, 0)
		//beego.Trace("Loading language: " + lang)
		filepath.Walk("resources/lang/"+lang, func(theFile string, f os.FileInfo, err error) error {
			//fmt.Println(theFile, f.Name(), err)
			if f != nil && !f.IsDir() {
				fileSuffix := path.Ext(f.Name())
				//fmt.Println("fileSuffix:",fileSuffix)
				if fileSuffix == ".ini" {
					//fmt.Println("ini:",theFile, f.Name())
					tempData, e := ioutil.ReadFile(theFile)
					if e != nil {
						beego.Error("Fail to set message file: " + err.Error())
						return err
					}
					langData = append(langData, tempData...)
				}
			}
			//fmt.Println("langData:", string(langData))
			return nil
		})
		if err := i18n.SetMessageData(lang, langData); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

// change orm db error to customer error
func DbError(err error) error {
	var result string
	switch err {
	case orm.ErrNoRows:
		result = "common.ormErrNoRows"
	case orm.ErrMissPK:
		result = "common.ormErrMissPK"
	case orm.ErrTxHasBegan:
		result = "common.ormErrTxHasBegan"
	case orm.ErrTxDone:
		result = "common.ormErrTxDone"
	case orm.ErrMultiRows:
		result = "common.ormErrMultiRows"
	case orm.ErrStmtClosed:
		result = "common.ormErrStmtClosed"
	case orm.ErrArgs:
		result = "common.ormErrArgs"
	case orm.ErrNotImplement:
		result = "common.ormErrNotImplement"
	default:
		result = err.Error()
	}
	return errors.New(result)
}

// change golang secret hash error to customer error
func HashError(err error) error {
	var result string
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		result = "common.hashErrMismatchedHashAndPassword"
	case bcrypt.ErrHashTooShort:
		result = "common.hashErrHashTooShort"
	default:
		result = err.Error()
	}
	return errors.New(result)
}

// change datetime to customer format
//  "2006-01-02 15:04:05.999999999 -0700 MST"
func DateFormat(datetime time.Time, format string) string {
	//fmt.Println("format date info :",datetime, format)
	defaultFormat := "2006-01-02 15:04:05"
	if format == "" {
		return datetime.Format(defaultFormat)
	}
	format = strings.Replace(format, "s", "05", -1)   //second
	format = strings.Replace(format, "i", "04", -1)   //minute
	format = strings.Replace(format, "H", "15", -1)   //hour
	format = strings.Replace(format, "d", "02", -1)   //day
	format = strings.Replace(format, "j", "02", -1)   //day
	format = strings.Replace(format, "m", "01", -1)   //month
	format = strings.Replace(format, "n", "01", -1)   //month
	format = strings.Replace(format, "Y", "2006", -1) //year
	return datetime.Format(format)
}

// format blog url
func WordPressUrlFormat(datetime time.Time, slug string, blogId int64, format string) string {
	//fmt.Println("WordPressUrlFormat date info :",datetime, slug, blogId, format)
	if format == "" {
		return "/?p=" + string(strconv.FormatInt(blogId, 10))
	}
	if strings.Contains(format, "archives") {
		return "/archives/" + string(strconv.FormatInt(blogId, 10))
	}
	format = strings.Replace(format, "%second%", "05", -1)   //second
	format = strings.Replace(format, "%minute%", "04", -1)   //moment
	format = strings.Replace(format, "%hour%", "15", -1)     //hour
	format = strings.Replace(format, "%day%", "02", -1)      //day
	format = strings.Replace(format, "%monthnum%", "01", -1) //month
	format = strings.Replace(format, "%year%", "2006", -1)   //year
	tmpFormat := datetime.Format(format)
	tmpFormat = strings.Replace(tmpFormat, "%post_id%", string(strconv.FormatInt(blogId, 10)), -1)
	tmpFormat = strings.Replace(tmpFormat, "%postname%", slug, -1)
	return tmpFormat
}

// page
func PageList(count, pageNo, pageSize int, url string, showNum int, lang string) string {
	param := "page="
	pageInfo := param + strconv.Itoa(pageNo)
	if strings.Contains(url, pageInfo) {
		url = strings.Replace(url, "?"+pageInfo, "", -1)
		url = strings.Replace(url, "&"+pageInfo, "", -1)
	}
	connectSign := "?"
	if strings.Contains(url, connectSign) {
		connectSign = "&"
	}
	if lang == "" {
		lang = "zh-CN"
	}
	if showNum == 0 {
		showNum = 3
	}
	totalPage := count / pageSize
	if count%pageSize > 0 {
		totalPage = count/pageSize + 1
	}

	firstPageString := `<li class="page-item"><a class="page-link" href="` + url + connectSign + param + "1" + `">` + i18n.Tr(lang, "common.firstPage") + `</a></li>`
	if pageNo <= 1 {
		//firstPageString = `<li class="page-item disabled"><a class="page-link" tabindex="-1" href="` + url + connectSign + param + "1" + `">` + i18n.Tr(lang, "common.firstPage") + `</a></li>`
		firstPageString = ""
	}
	lastPageString := `<li class="page-item"><a class="page-link" href="` + url + connectSign + param + strconv.Itoa(totalPage) + `">` + i18n.Tr(lang, "common.lastPage") + `</a></li>`
	if pageNo >= totalPage {
		//lastPageString = `<li class="page-item disabled"><a class="page-link" tabindex="-1" href="` + url + connectSign + param + strconv.Itoa(totalPage) + `">` + i18n.Tr(lang, "common.lastPage") + `</a></li>`
		lastPageString = ""
	}

	previousPageString := `<li class="page-item"><a class="page-link" href="` + url + connectSign + param + strconv.Itoa(pageNo-1) + `">` + i18n.Tr(lang, "common.previousPage") + `</a></li>`
	if pageNo <= 2 {
		previousPageString = ""
	}
	nextPageString := `<li class="page-item"><a class="page-link" href="` + url + connectSign + param + strconv.Itoa(pageNo+1) + `">` + i18n.Tr(lang, "common.nextPage") + `</a></li>`
	if pageNo >= totalPage-2 {
		nextPageString = ""
	}

	prePageString := ""
	startPageSlid := pageNo - showNum
	if startPageSlid < 1 {
		startPageSlid = 1
	} else if startPageSlid > 1 {
		prePageString = `<li class="page-item disabled"><span class="page-link">...</span></li>`
	}
	for i := startPageSlid; i < pageNo; i++ {
		prePageString += `<li class="page-item"><a class="page-link" href="` + url + connectSign + param + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
	}

	subPageString := ""
	endPageSlid := pageNo + showNum
	if endPageSlid > totalPage {
		endPageSlid = totalPage
	} else if endPageSlid < totalPage {
		subPageString = `<li class="page-item disabled"><span class="page-link">...</span></li>`
	}
	for i := endPageSlid; i > pageNo; i-- {
		subPageString = `<li class="page-item"><a class="page-link" href="` + url + connectSign + param + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>` + subPageString
	}
	currentPageString := `<li class="page-item active"><span class="page-link">` + strconv.Itoa(pageNo) + `<span class="sr-only">(current)</span></span></li>`
	pageList := firstPageString + previousPageString + prePageString + currentPageString + subPageString + nextPageString + lastPageString
	return `<nav class="page" aria-label="` + i18n.Tr(lang, "common.pageAriaLabel") + `"><ul class="pagination justify-content-center">` + pageList + `</ul></nav>`
}
