package pkg

import (
	"fmt"
	"testing"

	"github.com/everfore/exc"
	"github.com/toukii/goutils"
)

func TestPkg(t *testing.T) {
	bs, err := exc.NewCMD("go list -json").DoNoTime()
	gopkg, err := NewPkg(bs)
	if goutils.CheckErr(err) {
		return
	}
	fmt.Println(gopkg.DepPkgsIndent())
	fmt.Println(gopkg.ImportPkgsIndent())
}
