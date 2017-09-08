package exc

import (
	"fmt"
	"testing"
	"time"
)

func TestExcAfter(t *testing.T) {
	NewCMD("go version").Debug().ExecuteAfter(4)
}

func TestExecHere2(t *testing.T) {
	var retbs []byte
	bs, err := NewCMD("go env").Debug().Env("GOPATH").Execute().Cd("src").Reset("go version").Out(&retbs, nil).Do()
	fmt.Println(string(bs))
	fmt.Println(err)

	fmt.Println(string(retbs))

	NewCMD("git push origin master").Debug().Env("GOPATH").Cd("src/github.com/everfore/exc").Wd().Execute()
}

func TestExecHere3(t *testing.T) {
	bs, err := NewCMD("ping www.baidu.com").Debug().Env("GOPATH").Do() //.Cd("src").Reset("go version").Out(&retbs, nil).Do()
	fmt.Println(string(bs), err)
	time.Sleep(8e9)
}

func TestExecHere4(t *testing.T) {
	bs, _ := NewCMD("go -L").Env("GOPATH").Do()
	fmt.Println(string(bs))
	// Checkerr(err)
}

func TestDoNoTime(t *testing.T) {
	NewCMD("go version").Debug().DoNoTime()
}

func TestBash(t *testing.T) {
	bs, err := Bash("cat exc.go | grep cmd | wc -l").DoNoTime()
	t.Log(string(bs), err)
	Bash("cat exc.go | grep cmd | wc").Debug().ExecuteAfter(1).Out(&bs, &err)
	t.Log(string(bs), err)
}
