package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/toukii/goutils"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	tplFile  = kingpin.Arg("tpl", "template file").Default("s2b.tpl").String()
	dataFile = kingpin.Arg("data", "render data file").Default("data.yaml").String()

	// gob = kingpin.Flag("pkg", "go build pkg").Short('p').String()
	// status          = kingpin.Command("status", "Check lock status")
	// release         = kingpin.Command("release", "Release the lock, uses .consul_lock_id file")
	// release_with_id = kingpin.Command("release-with-id", "Release the lock with explicit session ID argument. Should not be used normally")
	// releaseID       = release_with_id.Arg("id", "session ID to release").Required().String()
)

type Data struct {
	Map map[string]interface{} `yaml:"map"`
	Arr []interface{}          `yaml:"arr"`
}

func Test(arg int) string {
	return fmt.Sprintf("TEST call %d", arg)
}

func main() {
	_ = kingpin.Parse()

	if *tplFile == "" {
		fmt.Println("template file is nil.")
		os.Exit(-1)
	}

	dfbs := goutils.ReadFile(*dataFile)
	var data Data
	err := yaml.Unmarshal(dfbs, &data)
	if goutils.LogCheckErr(err) {
		os.Exit(-1)
	}

	err = Render(*tplFile, data)
	goutils.LogCheckErr(err)
	// if len(data.Map) > 0 {
	// 	err := Render(*tplFile, data.Map)
	// 	goutils.LogCheckErr(err)
	// }
	// if len(data.Arr) > 0 {
	// 	err := Render(*tplFile, data.Arr)
	// 	goutils.LogCheckErr(err)
	// }

}

func Render(tplFile string, data interface{}) error {
	tpl := template.New(tplFile).Funcs(template.FuncMap{"test": func(arg int) string {
		return fmt.Sprintf("TEST call %d", arg)
	}})
	tpl, err := tpl.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	err = tpl.ExecuteTemplate(os.Stdout, tplFile, data)
	return err
}
