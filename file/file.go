
package kweb_util_file

import (
	os "os"
	io "io"
)

func IsExist(path string)(bool){
	s, err := os.Stat(path)
	return (s != nil) || (err != nil && os.IsExist(err))
}

func IsNotExist(path string)(bool){
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func IsFile(path string)(bool){
	s, _ := os.Stat(path)
	return s != nil && !s.IsDir()
}

func IsDir(path string)(bool){
	s, _ := os.Stat(path)
	return s != nil && s.IsDir()
}

func RemoveFile(path string)(bool){
	err := os.Remove(path)
	return err == nil
}

func CreateDir(folderPath string, mode_ ...os.FileMode)(err error){
	mode := os.ModePerm
	if len(mode_) > 0 {
		mode = mode_[0]
	}
	if IsNotExist(folderPath){
		err = os.Mkdir(folderPath, os.ModePerm)
		if err != nil { return }
		err = os.Chmod(folderPath, mode)
		if err != nil { return }
	}
	return nil
}

func CopyFile(src string, drt string)(err error){
	var (
		sfd *os.File
		dfd *os.File
		info os.FileInfo
	)
	sfd, err = os.Open(src)
	if err != nil { return }
	defer sfd.Close()

	dfd, err = os.Open(drt)
	if err != nil { return }
	defer dfd.Close()

	_, err = io.Copy(sfd, dfd)
	if err != nil { return }

	info, err = sfd.Stat()
	if err != nil { return }

	err = dfd.Chmod(info.Mode())
	return
}

func CopyDir(src string, drt string)(err error){
	var (
		sfd *os.File
		dirinfo os.FileInfo
		finfos []os.FileInfo
	)
	sfd, err = os.Open(src)
	if err != nil { return }
	defer sfd.Close()

	dirinfo, err = sfd.Stat()
	if err != nil { return }

	err = CreateDir(drt, dirinfo.Mode())
	if err != nil { return }
	defer func(){ if err != nil {
		os.RemoveAll(drt)
	}}()

	finfos, err = sfd.Readdir(-1)
	if err != nil { return }

	sfd.Close()

	for _, info := range finfos {
		sf := JoinPath(src, info.Name())
		df := JoinPath(drt, info.Name())
		if info.IsDir() {
			err = CopyDir(df, df)
		}else{
			err = CopyFile(sf, df)
		}
		if err != nil { return }
	}
	return nil
}

