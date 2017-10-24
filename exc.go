package exc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

type CMD struct {
	raw       string
	cmd       string
	args      []string
	debug     bool
	Execution func(cmd string, args ...string) ([]byte, error)
}

func NewCMD(cmd string) *CMD {
	newCMD := &CMD{raw: cmd, Execution: DefaultExecution}
	args := strings.Split(cmd, " ")
	args_len := len(args)
	if args_len > 1 {
		newCMD.cmd = args[0]
		newCMD.args = args[1:]
	} else {
		newCMD.cmd = cmd
	}
	return newCMD
}

func Bash(cmd string) *CMD {
	return &CMD{raw: cmd, cmd: cmd, Execution: BashExecution}
}

func (c *CMD) Debug() *CMD {
	c.debug = !c.debug
	return c
}

func (c *CMD) Cd(dir string) *CMD {
	err := os.Chdir(dir)
	if Checkerr(err) {
		panic(err)
	}
	return c
}

func (c *CMD) Env(key string) *CMD {
	return c.Cd(os.Getenv(key))
}

func (c *CMD) Wd() *CMD {
	dir, err := os.Getwd()
	if !Checkerr(err) {
		fmt.Printf("=======[ %s ]======\n", dir)
	}
	return c
}

func (c *CMD) Reset(cmd string) *CMD {
	c.raw = cmd
	args := strings.Split(cmd, " ")
	args_len := len(args)
	if args_len > 1 {
		c.cmd = args[0]
		c.args = args[1:]
	} else {
		c.cmd = cmd
	}
	return c
}

func (c *CMD) DoNoTime() ([]byte, error) {
	if c.debug {
		fmt.Printf("\n## cmd:\n%s", c.raw)
	}
	if len(c.raw) <= 0 {
		return nil, fmt.Errorf("raw cmd is nil")
	}
	if c.Execution != nil {
		return c.Execution(c.cmd, c.args...)
	}
	return nil, nil
}

func (c *CMD) Do() ([]byte, error) {
	if c.debug {
		fmt.Printf("\n## cmd:\n%s", c.raw)
	}
	if len(c.raw) <= 0 {
		return nil, fmt.Errorf("raw cmd is nil")
	}
	var err error
	if c.Execution != nil {
		var bs []byte
		sh := make(chan bool)
		go func() {
			go func() {
				<-time.After(8e9)
				runtime.Goexit()
			}()
			bs, err = c.Execution(c.cmd, c.args...)
			sh <- true
		}()
		select {
		case <-sh:
			Checkerr(err)
			if c.debug {
				fmt.Println(*(*string)(unsafe.Pointer(&bs)))
			}
			return bs, err
		case <-time.After(8e9):
			err = fmt.Errorf("Timeout")
		}
	}
	return nil, err
}

func (c *CMD) Execute() *CMD {
	bs, err := c.DoNoTime()
	if err != nil {
		fmt.Println("\n## execute err:", err)
	}
	if c.debug {
		fmt.Printf("\n## execute output:\n%s", *(*string)(unsafe.Pointer(&bs)))
	}
	return c
}

func (c *CMD) ExecuteAfter(sec int) *CMD {
	<-time.After(time.Second * (time.Duration)(sec))
	c.Do()
	return c
}

func (c *CMD) Out(ret *[]byte, reterr *error) *CMD {
	b, err := c.Do()
	if ret != nil {
		*ret = b
	}
	if reterr != nil {
		*reterr = err
	}
	return c
}

func (c *CMD) Tmp(cmd string) *CMD {
	fmt.Printf("[Tmp cmd] %s\n", cmd)
	var excmd *exec.Cmd
	args := strings.Split(cmd, " ")
	args_len := len(args)
	if args_len > 1 {
		excmd = exec.Command(args[0], args[1:]...)
	} else {
		excmd = exec.Command(cmd)
	}
	bs, _ := excmd.CombinedOutput()
	fmt.Println(*(*string)(unsafe.Pointer(&bs)))
	return c
}

func DefaultExecution(cmd string, args ...string) ([]byte, error) {
	var excmd *exec.Cmd
	if len(args) > 0 {
		excmd = exec.Command(cmd, args...)
	} else {
		excmd = exec.Command(cmd)
	}
	return excmd.CombinedOutput()
}

func BashExecution(bash string, args ...string) ([]byte, error) {
	return exec.Command("bash", "-c", bash).CombinedOutput()
}

func Checkerr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
