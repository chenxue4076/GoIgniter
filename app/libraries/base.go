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
)


//var langs = []string {"zh-CN", "en-US"}

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