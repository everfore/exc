package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/everfore/exc"
	"github.com/everfore/exc/walkexc"
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
		size := len(args)
		if size > 0 {
			viper.Set("pkg", args[0])
		} else {
			viper.Set("pkg", ".")
		}
		walkexc.Setting(Cond, "")
		filepath.Walk(viper.GetString("pkg"), walkexc.WalkExc)
	},
}

var excuted map[string]bool

func init() {
	Command.PersistentFlags().BoolP("import", "i", true, "display import pkgs")
	Command.PersistentFlags().BoolP("dep", "d", false, "display dependency pkgs")
	Command.PersistentFlags().BoolP("error", "e", true, "display error pkgs")
	Command.PersistentFlags().BoolP("install", "I", false, "install error pkgs")
	Command.PersistentFlags().BoolP("R-repeat", "R", false, "R-repeat")
	viper.BindPFlag("import", Command.PersistentFlags().Lookup("import"))
	viper.BindPFlag("dep", Command.PersistentFlags().Lookup("dep"))
	viper.BindPFlag("error", Command.PersistentFlags().Lookup("error"))
	viper.BindPFlag("install", Command.PersistentFlags().Lookup("install"))
	viper.BindPFlag("R-repeat", Command.PersistentFlags().Lookup("R-repeat"))
	excuted = make(map[string]bool)
}

func hasGoDir(filename string) (string, bool) {
	if !strings.HasSuffix(filename, ".go") {
		return "", false
	}
	abs, _ := filepath.Abs(filename)
	return abs[:strings.LastIndex(abs, string(filepath.Separator))], true
}

func Cond(path string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.EqualFold(info.Name(), ".git") {
		return false, filepath.SkipDir
	}
	if !info.IsDir() {
		absBase, has := hasGoDir(path)
		if !has || excuted[absBase] {
			return false, nil
		}
		excuted[absBase] = true
		Excute(TrimGopath(absBase))
		return false, nil
	} else {
		if !viper.GetBool("R-repeat") {
			if path != viper.GetString("pkg") {
				return false, filepath.SkipDir
			}
		}
	}
	return false, nil
}

func Excute(repo string) error {
	bs, err := exc.NewCMD("go list -json " + repo).DoNoTime()
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
