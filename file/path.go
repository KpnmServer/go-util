
package kweb_util_file

import (
	os "os"
	strings "strings"
)

func JoinPath(paths ...string)(allpath string){
	allpath = ""
	for _, p := range paths {
		if len(p) == 0 {
			continue
		}
		if p[0] == '/' {
			allpath = p
			continue
		}
		if len(allpath) != 0 && allpath[len(allpath) - 1] != '/' {
			allpath += "/"
		}
		allpath += p
	}
	return allpath
}

func JoinPathWithoutAbs(paths ...string)(allpath string){
	allpath = ""
	for _, p := range paths {
		if len(p) == 0 {
			continue
		}
		if len(allpath) != 0 && allpath[len(allpath) - 1] != '/' && p[0] != '/' {
			allpath += "/"
		}
		allpath += p
	}
	return allpath
}

func SplitPaths(path string)(paths []string){
	if len(path) == 0 {
		return []string{}
	}
	paths = strings.Split(path, "/")
	if path[0] == '/' {
		paths[0] = "/"
	}
	if path[len(path) - 1] == '/' {
		paths = paths[:len(paths) - 1]
	}
	return
}

func RunPath()(cwdPath string){
	var err error
	cwdPath, err = os.Getwd()
	if err != nil {
		panic(err)
		return "."
	}
	return cwdPath
}

func FixPath(path string)(string){
	path0 := SplitPaths(path)
	paths := make([]string, len(path0))
	var i int = 0
	for _, p := range path0 {
		if p == "." {
			continue
		}
		if p == ".." {
			i--
		}else{
			if i >= 0 {
				paths[i] = p
			}
			i++
		}
	}
	if i <= 0 {
		if path[0] == '/' {
			return "/"
		}
		if i == 0 {
			return "."
		}
		paths = make([]string, -i)
		for n, _ := range paths {
			paths[n] = ".."
		}
		return JoinPath(paths...)
	}
	return JoinPath(paths[:i]...)
}

func AbsPath(path string)(string){
	if path[0] == '/' {
		return path
	}
	return JoinPath(RunPath(), path)
}

func RelPath(path string, base_ ...string)(string){
	var bases []string
	if len(base_) > 0 {
		bases = SplitPaths(FixPath(AbsPath(base_[0])))
	}else{
		bases = SplitPaths(RunPath())
	}
	paths := SplitPaths(FixPath(AbsPath(path)))
	var mlen = len(bases)
	if mlen > len(paths) { mlen = len(paths) }
	var x int = 0
	for ; x < mlen ;x++ {
		if paths[x] != bases[x] {
			break
		}
	}
	var backs []string = make([]string, len(bases) - x)
	for i, _ := range backs {
		backs[i] = ".."
	}
	return JoinPath(append(backs, paths[x:]...)...)
}


func SplitPath(path string)(dirn string, base string){
	var i int
	for i = len(path) - 1; i >= 0 && path[i] != '/' ;i-- {}
	if i == -1 {
		return "", path
	}
	return path[:i], path[i + 1:]
}

func DirPath(path string)(dirn string){
	dirn, _ = SplitPath(path)
	return
}

func BasePath(path string)(base string){
	_, base = SplitPath(path)
	return
}

func SplitName(path string)(base string, suffix string){
	var i int
	for i = len(path) - 1; i >= 0 && path[i] != '.' ;i-- {}
	if i == -1 {
		return path, ""
	}
	return path[:i], path[i + 1:]
}

func BaseName(path string)(base string){
	base, _ = SplitName(path)
	return
}

func SuffixName(path string)(suffix string){
	_, suffix = SplitName(path)
	return
}
