package walkexc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	command  []string
	locker   sync.Mutex
	CondFunc func(path string, info os.FileInfo) (ignore bool, skip error)
)

func init() {
	CondFunc = DefaultCond
}

func Setting(cond func(path string, info os.FileInfo) (ignore bool, skip error), c ...string) {
	locker.Lock()
	defer locker.Unlock()
	command = c
	if cond == nil {
		CondFunc = DefaultCond
	} else {
		CondFunc = cond
	}
}

func setCmd(c ...string) *exec.Cmd {
	locker.Lock()
	defer locker.Unlock()
	leng := len(c)
	var cmd *exec.Cmd
	if leng < 1 {
		cmd = exec.Command("ls")
	} else if leng == 1 {
		cmd = exec.Command(c[0])
	} else {
		cmd = exec.Command(c[0], c[1:]...)
	}
	return cmd
}

func DefaultCond(path string, info os.FileInfo) (ignore bool, skip error) {
	if strings.EqualFold(info.Name(), ".git") {
		return true, filepath.SkipDir
	}
	if !info.IsDir() {
		return true, nil
	}
	return false, nil
}

func WalkExc(path string, info os.FileInfo, err error) error {
	if ignore, skip := CondFunc(path, info); skip != nil || ignore {
		return skip
	}
	abs, _ := filepath.Abs(path)
	fmt.Println("cur path:", abs)
	cmd := setCmd(command...)
	cmd.Dir = filepath.Clean(path)
	bs, _ := cmd.CombinedOutput()
	fmt.Println(string(bs))
	return nil
}
