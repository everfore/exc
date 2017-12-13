package rep

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/everfore/exc/walkexc"
	cr "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toukii/goutils"
)

var Command = &cobra.Command{
	Use:   "rep",
	Short: "replace file content",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("old") == "" && viper.GetString("new") == "" {
			cmd.Help()
			return
		}
		newBytes = goutils.ToByte(viper.GetString("new"))
		walkexc.Setting(Cond, "")
		filepath.Walk(viper.GetString("dir"), walkexc.WalkExc)
	},
}

var newBytes []byte

func init() {
	Command.PersistentFlags().StringP("dir", "p", ".", "the root dir of replace files")
	Command.PersistentFlags().StringP("skip", "s", ".git", "skip dir")
	Command.PersistentFlags().BoolP("repeat", "r", false, "repeat to execute")
	Command.PersistentFlags().StringP("extension", "e", ".go", "file extension, ex: .go ; Dockerfilw")
	Command.PersistentFlags().StringP("old", "o", "", "old content")
	Command.PersistentFlags().StringP("new", "n", "", "new content")
	viper.BindPFlag("dir", Command.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("skip", Command.PersistentFlags().Lookup("skip"))
	viper.BindPFlag("repeat", Command.PersistentFlags().Lookup("repeat"))
	viper.BindPFlag("extension", Command.PersistentFlags().Lookup("extension"))
	viper.BindPFlag("old", Command.PersistentFlags().Lookup("old"))
	viper.BindPFlag("new", Command.PersistentFlags().Lookup("new"))
}

func Cond(path string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.EqualFold(info.Name(), viper.GetString("skip")) {
		return false, filepath.SkipDir
	}
	if !info.IsDir() {
		if strings.HasSuffix(info.Name(), viper.GetString("extension")) {
			if bs, cnt := ContentContains(path, viper.GetString("old")); cnt {
				fmt.Printf("rewrite: %s\n", cr.HiMagentaString(path))
				rst_cnt := regexp.MustCompile(viper.GetString("old")).ReplaceAllFunc(bs, func(b []byte) []byte {
					return newBytes
				})
				goutils.ReWriteFile(path, rst_cnt)
				return true, nil
			}
		}
	} else {
		if !viper.GetBool("repeat") && path != viper.GetString("dir") {
			fmt.Println("skip path:", path)
			return false, filepath.SkipDir
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
