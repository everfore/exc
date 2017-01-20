package walkexc

import (
	// "fmt"
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

func DefaultCond(path string, info os.FileInfo) (ifExec bool, skip error) {
	if strings.EqualFold(info.Name(), ".git") {
		return false, filepath.SkipDir
	}
	if !info.IsDir() {
		return true, nil
	}
	return true, nil
}

func WalkExc(path string, info os.FileInfo, err error) error {
	ifExec, skip := CondFunc(path, info)
	if skip != nil {
		return skip
	}
	if ifExec {
		if info.IsDir() {
			PathExc(path, command...)
		} else {
			base_path := filepath.Dir(path)
			PathExc(base_path, command...)
		}
	}
	return nil
}

func PathExc(path string, c ...string) error {
	abs, _ := filepath.Abs(path)
	// fmt.Println("cur path:", abs, c)
	cmd := setCmd(c...)
	cmd.Dir = abs
	// bs, _ := cmd.CombinedOutput()
	// fmt.Println(string(bs))
	cmd.CombinedOutput()
	return nil
}
