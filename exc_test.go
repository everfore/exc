package exc

import (
	"fmt"
	"testing"
)

func TestExecHere2(t *testing.T) {
	var retbs []byte
	bs, err := NewCMD("go env").Debug().Env("GOPATH").Execute().Cd("src").Reset("go version").Out(&retbs, nil).Do()
	fmt.Println(string(bs))
	fmt.Println(err)

	fmt.Println(string(retbs))

}
