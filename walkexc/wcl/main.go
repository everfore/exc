package main

import (
	"bufio"
	"fmt"
	"github.com/everfore/exc/walkexc"
	"github.com/toukii/goutils"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	count     = 0
	fileCount = 0
	mc        map[string]int
)

func init() {
	mc = make(map[string]int)
	mc[".go"] = 0
}

func main() {
	walkexc.Setting(Cond, "")
	filepath.Walk("./", walkexc.WalkExc)
	fmt.Println("wcl:", count)
	fmt.Println(mc)
}

func Cond(path1 string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.Contains(path1, ".git/") || strings.Contains(path1, "vendor/") ||
		strings.HasPrefix(info.Name(), "gen.") || strings.HasPrefix(info.Name(), "gen_") || strings.HasSuffix(info.Name(), ".pb.go") {
		return false, nil
	}
	fileSuffix := path.Ext(info.Name())
	if !info.IsDir() {
		if ".go" == fileSuffix {
			ct := Wcl(path1)
			mc[fileSuffix] += ct
			count += ct
			fileCount++
			// fmt.Printf("%d ", fileCount)
			return true, nil
		}
	}
	return false, nil
}

func Wcl(filaname string) int {
	abs_path, err := filepath.Abs(filaname)
	if goutils.CheckErr(err) {
		return 0
	}
	fr, err := os.OpenFile(abs_path, os.O_RDONLY, 0644)
	defer fr.Close()
	if goutils.CheckErr(err) {
		return 0
	}
	br := bufio.NewReader(fr)
	ct := 0
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			break
		}
		if len(line) <= 0 {
			// continue
		}
		ct++
	}
	// fmt.Printf("%s:%d\n", abs_path, ct)
	return ct
}

func ContentContains(path string, stuff string) ([]byte, bool) {
	abs_path, err := filepath.Abs(path)
	if goutils.CheckErr(err) {
		return nil, false
	}
	fr, err := os.OpenFile(abs_path, os.O_RDONLY, 0644)
	defer fr.Close()
	if goutils.CheckErr(err) {
		return nil, false
	}
	bs, err := ioutil.ReadAll(fr)
	if goutils.CheckErr(err) {
		return nil, false
	}
	if strings.Contains(goutils.ToString(bs), stuff) {
		return bs, true
	}
	return nil, false
}
