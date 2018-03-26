package tpl

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/everfore/exc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	qiniubytes "github.com/toukii/bytes"
	"github.com/toukii/goutils"
	yaml "gopkg.in/yaml.v2"
)

var Command = &cobra.Command{
	Use:   "tpl",
	Short: "render template",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		Excute(args)
	},
}

func init() {
	Command.PersistentFlags().BoolP("execute", "e", false, "excute the render result by bash")
	viper.BindPFlag("execute", Command.PersistentFlags().Lookup("execute"))

	Command.AddCommand(InitCommand)
}

func Excute(args []string) error {
	newArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		newArgs = append(newArgs, arg)
	}
	size := len(newArgs)
	if size == 1 {
		viper.Set("tpl", newArgs[0])
	} else if size >= 2 {
		viper.Set("tpl", newArgs[0])
		viper.Set("data", newArgs[1])
	} else {
		viper.Set("tpl", "tpl.tpl")
		viper.Set("data", "data.yaml")
	}

	return excute(viper.GetString("tpl"), viper.GetString("data"), viper.GetBool("execute"))
}

// Data tpl data
type Data struct {
	Map  map[string]interface{} `yaml:"map"`
	Arr  []interface{}          `yaml:"arr"`
	Step int
}

func excute(tplFile string, dataFile string, excute bool) error {
	if tplFile == "" {
		return fmt.Errorf("template file is nil.")
	}

	dfbs := goutils.ReadFile(dataFile)
	var data Data
	err := yaml.Unmarshal(dfbs, &data)
	if goutils.CheckNoLogErr(err) {
		return err
	}

	if excute {
		buf := make([]byte, 20, 10240)
		wr := qiniubytes.NewWriter(buf)

		if data.Step > 0 {
			size := len(data.Arr)
			N := data.Step
			remainder := 0 // 没有没整除的部分
			if size%N > 0 {
				remainder = 1
			}
			loop := size/N + remainder // 总共切片切次
			for i := 0; i < loop-1; i++ {
				dt := &Data{
					Map: data.Map,
					Arr: data.Arr[i*N : (i+1)*N],
				}
				wr.Reset()
				err := Render(wr, tplFile, dt)
				if goutils.CheckNoLogErr(err) {
					return err
				}
				exc.Bash(goutils.ToString(wr.Bytes())).Debug().Execute()
			}
			dt := &Data{
				Map: data.Map,
				Arr: data.Arr[(loop-1)*N:],
			}
			err := Render(wr, tplFile, dt)
			if goutils.CheckNoLogErr(err) {
				return err
			}
			exc.Bash(goutils.ToString(wr.Bytes())).Debug().Execute()
		} else {
			err := Render(wr, tplFile, data)
			if goutils.CheckNoLogErr(err) {
				return err
			}
			exc.Bash(goutils.ToString(wr.Bytes())).Debug().Execute()
		}
	} else {
		return Render(os.Stdout, tplFile, data)
	}
	return nil
}

// Render
func Render(wr io.Writer, tplFile string, data interface{}) error {
	tpl := template.New(tplFile).Funcs(template.FuncMap{
		"join": func(ids []interface{}, sep string) string {
			idsStr := make([]string, len(ids))
			for i, id := range ids {
				idsStr[i] = fmt.Sprint(id)
			}
			return strings.Join(idsStr, sep)
		},
		"loop": func(start, end, step int) []int {
			if step <= 0 || start > end {
				return nil
			}
			loop := make([]int, (end-start)/step+1)
			i := 0
			for start <= end {
				loop[i] = start
				start += step
				i++
			}
			return loop
		},
	})
	tpl, err := tpl.Parse(goutils.ToString(goutils.ReadFile(tplFile)))
	if goutils.CheckNoLogErr(err) {
		return err
	}

	return tpl.Execute(wr, data)
}
