package tpl

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toukii/goutils"
)

var (
	InitCommand = &cobra.Command{
		Use:   "init",
		Short: "init render template",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			InitExcute()
		},
	}

	tpl = `{{range $k,$v := .Map}}
	echo hello, {{$k}}:{{$v}}
	{{range $v}}
		echo name: {{.First}} {{.Second}}
	{{end}}
{{end}}

{{range $i,$v := .Arr}}
	echo arr[{{$i}}]:{{$v}}
{{end}}

{{with .Arr}}
	echo arrs: {{join . ","}}___
{{end}}

{{loop 1 10 1}}
{{loop 1 10 2}}
{{loop 1 10 3}}

{{range $i:=loop 1 10 2}}
	{{$i}}th
{{end}}
`

	data = `map:
  Names: 
    - 
      First: Simth
      Second: John
    - 
      First: Simth
      Second: John2
    - 
      First: Simth
      Second: John3
arr:
  - 1
  - 2
  - 3
`
)

func init() {
	InitCommand.PersistentFlags().StringP("output", "o", ".", "output dir")
	viper.BindPFlag("output", InitCommand.PersistentFlags().Lookup("output"))
}

func InitExcute() error {
	goutils.WriteFile(filepath.Join(viper.GetString("output"), "tpl.tpl"), goutils.ToByte(tpl))
	return goutils.WriteFile(filepath.Join(viper.GetString("output"), "data.yaml"), goutils.ToByte(data))
}
