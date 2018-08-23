package libraries

import (
	"path/filepath"
	"os"
	"path"
	"io/ioutil"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"fmt"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"time"
	"strings"
	"strconv"
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
func LoadLangs()  {
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
			if f != nil && ! f.IsDir() {
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
	format = strings.Replace(format,"s", "05", -1)	//second
	format = strings.Replace(format,"i", "04", -1)	//minute
	format = strings.Replace(format,"H", "15", -1)	//hour
	format = strings.Replace(format,"d", "02", -1)	//day
	format = strings.Replace(format,"j", "02", -1)	//day
	format = strings.Replace(format,"m", "01", -1)	//month
	format = strings.Replace(format,"n", "01", -1)	//month
	format = strings.Replace(format,"Y", "2006", -1)	//year
	return datetime.Format(format)
}
// format blog url
func WordPressUrlFormat(datetime time.Time, slug string, blogId int64, format string ) string {
	//fmt.Println("WordPressUrlFormat date info :",datetime, slug, blogId, format)
	if format == "" {
		return "/?p="+ string(strconv.FormatInt(blogId, 10))
	}
	if strings.Contains(format, "archives") {
		return "/archives/" + string(strconv.FormatInt(blogId, 10))
	}
	format = strings.Replace(format,"%second%", "05", -1)	//second
	format = strings.Replace(format,"%minute%", "04", -1)	//moment
	format = strings.Replace(format,"%hour%", "15", -1)	//hour
	format = strings.Replace(format,"%day%", "02", -1)	//day
	format = strings.Replace(format,"%monthnum%", "01", -1)	//month
	format = strings.Replace(format,"%year%", "2006", -1)	//year
	tmpFormat := datetime.Format(format)
	tmpFormat = strings.Replace(tmpFormat,"%post_id%", string(strconv.FormatInt(blogId, 10)), -1)
	tmpFormat = strings.Replace(tmpFormat,"%postname%", slug, -1)
	return tmpFormat
}

