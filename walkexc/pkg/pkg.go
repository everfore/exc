package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cr "github.com/fatih/color"
	"github.com/toukii/goutils"
)

type Pkg struct {
	Dir         string   `json:"Dir"`
	ImportPath  string   `json:"ImportPath"`
	Name        string   `json:"Name"`
	Target      string   `json:"Target"`
	Stale       bool     `json:"Stale"`
	StaleReason string   `json:"StaleReason"`
	Root        string   `json:"Root"`
	GoFiles     []string `json:"GoFiles"`
	Imports     []string `json:"Imports"`
	Deps        []string `json:"Deps"`
	Incomplete  bool     `json:"Incomplete"`
	DepsErrors  []struct {
		ImportStack []string `json:"ImportStack"`
		Pos         string   `json:"Pos"`
		Err         string   `json:"Err"`
	} `json:"DepsErrors"`
	TestGoFiles []string `json:"TestGoFiles"`
	TestImports []string `json:"TestImports"`

	NoStdImportOKPkgs []string
	NoStdDepOKPkgs    []string
	NoStdDepErrPkgs   []*NoStdDepErrPkg
	nostdDepErrPkgMap map[string]bool
}

type NoStdDepErrPkg struct {
	PkgName string
	Tips    string
}

func (ep *NoStdDepErrPkg) String() string {
	return cr.RedString("%s", ep.PkgName) + cr.HiRedString(" [%s]", ep.Tips)
}

func NewPkg(bs []byte) (*Pkg, error) {
	var pkg Pkg
	err := json.Unmarshal(bs, &pkg)
	if goutils.CheckNoLogErr(err) {
		return nil, fmt.Errorf("%s", bs)
	}
	pkg.analyse()
	return &pkg, err
}

func (p *Pkg) analyse() {
	p.NoStdDepErrPkgs = p.analyseDepsErrors()

	p.NoStdImportOKPkgs = make([]string, 0, len(p.Imports)-len(p.DepsErrors))
	for _, imps := range p.Imports {
		if p.nostdDepErrPkgMap[imps] || !GithubPkgFilter(imps) {
			continue
		}
		p.NoStdImportOKPkgs = append(p.NoStdImportOKPkgs, imps)
	}

	p.NoStdDepOKPkgs = make([]string, 0, len(p.Deps)-len(p.DepsErrors))
	for _, imps := range p.Deps {
		if p.nostdDepErrPkgMap[imps] || !GithubPkgFilter(imps) {
			continue
		}
		p.NoStdDepOKPkgs = append(p.NoStdDepOKPkgs, imps)
	}
}

func (p *Pkg) analyseDepsErrors() []*NoStdDepErrPkg {
	desize := len(p.DepsErrors)
	if desize <= 0 {
		return nil
	}
	p.nostdDepErrPkgMap = make(map[string]bool)
	ret := make([]*NoStdDepErrPkg, 0, desize)
	for _, depsErr := range p.DepsErrors {
		tips := depsErr.Err
		if strings.Contains(depsErr.Err, "cannot find package") {
			tips = "not-find"
		}
		for _, imp := range depsErr.ImportStack {
			if imp == p.ImportPath {
				continue
			}
			p.nostdDepErrPkgMap[imp] = true
			ret = append(ret, &NoStdDepErrPkg{
				PkgName: imp,
				Tips:    tips,
			})
		}
	}
	return ret
}

func (p *Pkg) ErrPkgsIndent() (string, bool) {
	return p.errPkgsIndent(1)
}

func (p *Pkg) errPkgsIndent(idx int) (indent string, output bool) {
	for _, name := range p.NoStdDepErrPkgs {
		output = true
		indent += fmt.Sprintf("\t%s. %s\n", cr.HiCyanString("%d", idx), name)
		idx++
	}
	return
}

func (p *Pkg) ImportPkgsIndent() (indent string, output bool) {
	idx := 1
	for _, name := range p.Imports {
		if GithubPkgFilter(name) {
			output = true
			indent += fmt.Sprintf("\t%s. %s\n", cr.HiCyanString("%d", idx), cr.HiGreenString(name))
			idx++
		}
	}

	if errIndent, errOutput := p.errPkgsIndent(idx); errOutput {
		indent += errIndent
	}
	return
}

func (p *Pkg) DepPkgsIndent() (indent string, output bool) {
	idx := 1
	for _, name := range p.NoStdDepOKPkgs {
		output = true
		indent += fmt.Sprintf("\t%s. %s\n", cr.HiCyanString("%d", idx), cr.HiGreenString(name))
		idx++
	}

	if errIndent, errOutput := p.errPkgsIndent(idx); errOutput {
		indent += errIndent
	}
	return
}

func TrimGopath(pkg string) string {
	abspath, err := filepath.Abs(os.Getenv("GOPATH"))
	if goutils.CheckErr(err) {
		return pkg
	}
	return strings.TrimPrefix(pkg, abspath+"/src/")
}

func GithubPkgFilter(pkg string) bool {
	if strings.Contains(pkg, ".com/") || strings.HasPrefix(pkg, "gopkg.in") || strings.HasPrefix(pkg, "golang.org") {
		return true
	}
	return false
}

/*{
	"Dir": "/Users/toukii/PATH/GOPATH/src/gitlab.1dmy.com/ezbuy/goflow/src/github.com/toukii/bkg",
	"ImportPath": "github.com/toukii/bkg",
	"Name": "main",
	"Target": "/Users/toukii/PATH/GOPATH/src/gitlab.1dmy.com/ezbuy/goflow/bin/bkg",
	"Stale": true,
	"StaleReason": "build ID mismatch",
	"Root": "/Users/toukii/PATH/GOPATH/src/gitlab.1dmy.com/ezbuy/goflow",
	"GoFiles": [
		"main.go"
	],
	"Imports": [
		"time"
	],
	"Deps": [
		"bufio",
		"unicode",
		"unicode/utf16",
		"unicode/utf8",
		"unsafe",
		"vendor/golang_org/x/crypto/chacha20poly1305",
		"vendor/golang_org/x/crypto/chacha20poly1305/internal/chacha20",
		"vendor/golang_org/x/crypto/curve25519",
		"vendor/golang_org/x/crypto/poly1305",
		"vendor/golang_org/x/net/http2/hpack",
		"vendor/golang_org/x/text/unicode/norm"
	],
	"Incomplete": true,
	"DepsErrors": [
		{
			"ImportStack": [
				"github.com/toukii/bkg",
				"github.com/everfore/exc"
			],
			"Pos": "main.go:11:2",
			"Err": "cannot find package \"github.com/everfore/exc\" in any of:\n\t/Users/toukii/PATH/Go/go/src/github.com/everfore/exc (from $GOROOT)\n\t/Users/toukii/PATH/GOPATH/src/gitlab.1dmy.com/ezbuy/goflow/src/github.com/everfore/exc (from $GOPATH)"
		}
	],
	"TestGoFiles": [
		"jsnm_test.go",
		"zhihu_test.go"
	],
	"TestImports": [
		"encoding/json",
		"fmt",
		"github.com/toukii/goutils",
		"github.com/toukii/membership/pkg3/go-simplejson",
		"io/ioutil",
		"net/http",
		"reflect",
		"testing"
	]
}
*/
