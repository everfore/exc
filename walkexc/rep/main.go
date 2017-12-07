package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/everfore/exc/walkexc"
	cr "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toukii/goutils"
)

var repcmd = cobra.Command{
	Use:   "rep",
	Short: "replace file content",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		walkexc.Setting(Cond, "")
		filepath.Walk(viper.GetString("dir"), walkexc.WalkExc)
	},
}

func init() {
	repcmd.PersistentFlags().StringP("dir", "p", ".", "the root dir of replace files")
	repcmd.PersistentFlags().StringP("skip", "s", ".git", "skip dir")
	repcmd.PersistentFlags().BoolP("repeat", "r", true, "repeat to execute")
	repcmd.PersistentFlags().StringP("extension", "e", ".go", "file extension, ex: .go ; Dockerfilw")
	repcmd.PersistentFlags().StringP("old", "o", "", "old content")
	repcmd.PersistentFlags().StringP("new", "n", "", "new content")
	viper.BindPFlag("dir", repcmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("skip", repcmd.PersistentFlags().Lookup("skip"))
	viper.BindPFlag("repeat", repcmd.PersistentFlags().Lookup("repeat"))
	viper.BindPFlag("extension", repcmd.PersistentFlags().Lookup("extension"))
	viper.BindPFlag("old", repcmd.PersistentFlags().Lookup("old"))
	viper.BindPFlag("new", repcmd.PersistentFlags().Lookup("new"))
}

func main() {
	if err := repcmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func Cond(path string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.EqualFold(info.Name(), viper.GetString("skip")) || !viper.GetBool("repeat") {
		return false, filepath.SkipDir
	}
	if !info.IsDir() {
		if strings.HasSuffix(info.Name(), viper.GetString("extension")) {
			if bs, cnt := ContentContains(path, viper.GetString("old")); cnt {
				fmt.Printf("rewrite: %s\n", cr.HiMagentaString(path))
				rst_cnt := strings.Replace(goutils.ToString(bs), viper.GetString("old"), viper.GetString("new"), -1)
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
	err := ioutil.WriteFile(path, bs, 0655)
	goutils.CheckErr(err)
}
