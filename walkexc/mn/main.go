package main

import (
	"fmt"
	"github.com/everfore/exc/walkexc"
	"github.com/toukii/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	// "bufio"
)

func main() {
	walkexc.Setting(Cond, "go", "version")
	filepath.Walk("./", walkexc.WalkExc)
}

func Cond(path string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.EqualFold(info.Name(), ".git") {
		return false, filepath.SkipDir
	}
	if !info.IsDir() {
		if strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "Dockerfile") {
			if bs, cnt := ContentContains(path, "toukii"); cnt {
				fmt.Println(info.Name())
				rst_cnt := strings.Replace(goutils.ToString(bs), "toukii", "toukii", -1)
				OverWrite(path, goutils.ToByte(rst_cnt))
				return true, nil
			}
		}
	}
	return false, nil
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

func OverWrite(path string, bs []byte) {
	// wf, err := os.OpenFile(path, os.O_WRONLY, 0655)
	// defer wf.Close()
	// if goutils.CheckErr(err) {
	// 	return
	// }
	err := ioutil.WriteFile(path, bs, 0655)
	goutils.CheckErr(err)
}
