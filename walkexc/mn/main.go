package main

import(
"github.com/everfore/exc/walkexc"
"path/filepath"
"fmt"
"os"
"strings"
)

func main() {
	walkexc.Setting(Cond, "go", "version")
	filepath.Walk("../", walkexc.WalkExc)
}

func Cond(path string, info os.FileInfo) (ignore bool, skip error){
	if strings.EqualFold(info.Name(), ".git") {
		return true, filepath.SkipDir
	}
	if !info.IsDir() {
		fmt.Println(info.Name())
		if strings.HasSuffix(info.Name(),".go") {
			
		}
		return true, nil
	}
	return false, nil
}