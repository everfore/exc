package exc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

func (c *CMD) Do() ([]byte, error) {
	if len(c.raw) <= 0 {
		return nil, fmt.Errorf("raw cmd is nil")
	}
	if c.Execution != nil {
		bs, err := c.Execution(c.cmd, c.args...)
		Checkerr(err)
		if c.debug {
			fmt.Println(*(*string)(unsafe.Pointer(&bs)))
		}
		return bs, err
	}
	return nil, nil
}

func (c *CMD) Execute() *CMD {
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

func DefaultExecution(cmd string, args ...string) ([]byte, error) {
	fmt.Printf("%s", cmd)
	var excmd *exec.Cmd
	if len(args) > 0 {
		excmd = exec.Command(cmd, args...)
		fmt.Printf(" %v\n", args)
	} else {
		excmd = exec.Command(cmd)
	}
	return excmd.CombinedOutput()
}

func Checkerr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
