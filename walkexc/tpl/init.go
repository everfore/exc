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
			if err := InitExcute(); goutils.CheckErr(err) {
			}
		},
	}
)

func init() {
	InitCommand.PersistentFlags().StringP("output", "o", ".", "output dir")
	viper.BindPFlag("output", InitCommand.PersistentFlags().Lookup("output"))
}

func InitExcute() error {
	if err := goutils.SafeWriteFile(filepath.Join(viper.GetString("output"), "tpl.tpl"), MustAsset("tpl/tpl.tpl")); goutils.CheckErr(err) {
	}
	return goutils.SafeWriteFile(filepath.Join(viper.GetString("output"), "data.yaml"), MustAsset("tpl/data.yaml"))
}
