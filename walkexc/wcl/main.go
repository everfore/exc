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
	count = 0
	mc    map[string]int
)

func init() {
	mc = make(map[string]int)
	mc[".java"] = 0
	mc[".js"] = 0
	mc[".jsp"] = 0
	mc[".html"] = 0
	mc[".properties"] = 0
}

func main() {
	walkexc.Setting(Cond, "")
	filepath.Walk("./", walkexc.WalkExc)
	fmt.Println("wcl:", count)
	fmt.Println(mc)
}

func Cond(path1 string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.HasPrefix(info.Name(), ".g") || strings.HasPrefix(info.Name(), ".s") || strings.HasPrefix(info.Name(), ".e") || strings.HasPrefix(info.Name(), "target") || strings.HasPrefix(info.Name(), "jquery") || strings.HasPrefix(info.Name(), "angular") || strings.HasSuffix(info.Name(), "min.js") {
		return false, filepath.SkipDir
	}
	if info.IsDir() && (strings.HasPrefix(info.Name(), "test") || strings.HasPrefix(info.Name(), "echarts") || strings.HasPrefix(info.Name(), "bootstrap") || strings.HasPrefix(info.Name(), "uikit") || strings.HasPrefix(info.Name(), "easyui") || strings.HasPrefix(info.Name(), "validate") || strings.HasPrefix(info.Name(), "highchart") || strings.HasPrefix(info.Name(), "textAngular") || strings.HasPrefix(info.Name(), "vendor")) {
		return false, filepath.SkipDir
	}
	fileSuffix := path.Ext(info.Name())
	if !info.IsDir() {
		if strings.Contains(".go.java.js.jsp.html.xml.properties", fileSuffix) {
			// if strings.Contains(".go.java", fileSuffix) {
			// if strings.Contains(".go.jsp", fileSuffix) {
			// if strings.Contains(".go.js", fileSuffix) {
			ct := Wcl(path1)
			mc[fileSuffix] += ct
			count += ct
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
	fmt.Printf("%s:%d\n", abs_path, ct)
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

func OverWrite(path string, bs []byte) {
	// wf, err := os.OpenFile(path, os.O_WRONLY, 0655)
	// defer wf.Close()
	// if goutils.CheckErr(err) {
	// 	return
	// }
	err := ioutil.WriteFile(path, bs, 0655)
	goutils.CheckErr(err)
}
