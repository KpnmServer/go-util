
package kpnmmail

import (
	ioutil "io/ioutil"

	ufile "github.com/KpnmServer/go-util/file"
)

var TEMPLATE_PATH = "emails"

func init(){
	{ // load mail template files
		var templateFiles []string = make([]string, 0)
		basePath := ufile.GetAbsPath(TEMPLATE_PATH)
		var findFunc func(path string)
		findFunc = func(path string){
			finfos, err := ioutil.ReadDir(path)
			if err != nil {
				panic(err)
			}
			for _, info := range finfos {
				fpath := ufile.JoinPath(path, info.Name())
				if info.IsDir() {
					findFunc(fpath)
				}else{
					templateFiles = append(templateFiles, fpath)
				}
			}
		}
		findFunc(basePath)

		if len(templateFiles) > 0 {
			LoadHtmlFiles(templateFiles...)
		}
	}
}
