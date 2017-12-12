package pkg

import (
	"fmt"

	"github.com/everfore/exc"
	cr "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toukii/goutils"
)

var Command = &cobra.Command{
	Use:   "pkg",
	Short: "go import/dependency packages",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		Excute()
	},
}

func init() {
	Command.PersistentFlags().BoolP("import", "i", true, "display import pkgs")
	Command.PersistentFlags().BoolP("dep", "d", false, "display dependency pkgs")
	Command.PersistentFlags().BoolP("error", "e", true, "display error pkgs")
	Command.PersistentFlags().BoolP("install", "I", false, "install error pkgs")
	viper.BindPFlag("import", Command.PersistentFlags().Lookup("import"))
	viper.BindPFlag("dep", Command.PersistentFlags().Lookup("dep"))
	viper.BindPFlag("error", Command.PersistentFlags().Lookup("error"))
	viper.BindPFlag("install", Command.PersistentFlags().Lookup("install"))
}

func Excute() error {
	bs, err := exc.NewCMD("go list -json").DoNoTime()
	gopkg, err := NewPkg(bs)
	if goutils.CheckErr(err) {
		return err
	}

	path := gopkg.ImportPath

	fmt.Printf("gopkg: [%s]\n", cr.HiCyanString(path))
	if viper.GetBool("dep") {
		indent, ok := gopkg.DepPkgsIndent()
		if ok {
			fmt.Printf("  %s:\n%s", cr.CyanString("dep gopkgs"), indent)
		}
	} else if viper.GetBool("import") {
		indent, ok := gopkg.ImportPkgsIndent()
		if ok {
			fmt.Printf("  %s:\n%s", cr.CyanString("import gopkgs"), indent)
		}
	}

	if viper.GetBool("error") {
		indent, ok := gopkg.ErrPkgsIndent()
		if ok {
			fmt.Printf("  %s:\n%s", cr.CyanString("err dep gopkgs"), indent)
		}
	}

	if viper.GetBool("install") {
		for _, pkg := range gopkg.NoStdDepErrPkgs {
			exc.NewCMD(fmt.Sprintf("pull %s", pkg.PkgName)).Debug().DoTimeout(20e9)
		}
	}
	return nil
}
