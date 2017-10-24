package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/everfore/exc"
	qiniubytes "github.com/toukii/bytes"
	"github.com/toukii/goutils"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	tplFile  = kingpin.Arg("tpl", "template file").Default("s2b.tpl").String()
	dataFile = kingpin.Arg("data", "render data file").Default("data.yaml").String()
	excute   = kingpin.Flag("excute", "excute command: the output").Short('E').Bool()

	// gob = kingpin.Flag("pkg", "go build pkg").Short('p').String()
	// status          = kingpin.Command("status", "Check lock status")
	// release         = kingpin.Command("release", "Release the lock, uses .consul_lock_id file")
	// release_with_id = kingpin.Command("release-with-id", "Release the lock with explicit session ID argument. Should not be used normally")
	// releaseID       = release_with_id.Arg("id", "session ID to release").Required().String()
)

// Data tpl data
type Data struct {
	Map map[string]interface{} `yaml:"map"`
	Arr []interface{}          `yaml:"arr"`
}

// Test
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

	if *excute {
		buf := make([]byte, 20, 10240)
		wr := qiniubytes.NewWriter(buf)
		err := Render(wr, *tplFile, data)
		if goutils.LogCheckErr(err) {
			os.Exit(-1)
		}
		// fmt.Fprintf(os.Stdout, "%s", wr.Bytes())
		exc.Bash(goutils.ToString(wr.Bytes())).Debug().Execute()
	} else {
		goutils.LogCheckErr(Render(os.Stdout, *tplFile, data))
	}
}

// Render
func Render(wr io.Writer, tplFile string, data interface{}) error {
	tpl := template.New(tplFile).Funcs(template.FuncMap{
		"test": func(arg int) string {
			return fmt.Sprintf("TEST call %d", arg)
		},
		"join": func(ids []interface{}, sep string) string {
			idsStr := make([]string, len(ids))
			for i, id := range ids {
				idsStr[i] = fmt.Sprint(id)
			}
			return strings.Join(idsStr, sep)
		},
	})
	tpl, err := tpl.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	return tpl.ExecuteTemplate(wr, tplFile, data)
}
