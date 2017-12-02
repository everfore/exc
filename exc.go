package exc

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	cr "github.com/fatih/color"
	"github.com/toukii/goutils"
)

type CMD struct {
	raw       string
	cmd       string
	args      []string
	debug     bool
	Execution func(cmd string, args ...string) ([]byte, error)
}

var (
	cmdfmt = cr.New(cr.FgCyan, cr.Bold, cr.BgYellow)
)

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
	if goutils.CheckNoLogErr(err) {
		panic(err)
	}
	return c
}

func (c *CMD) Env(key string) *CMD {
	return c.Cd(os.Getenv(key))
}

func (c *CMD) Wd() *CMD {
	dir, err := os.Getwd()
	if !goutils.CheckNoLogErr(err) {
		cr.Cyan("= = = [ %s ] = = =\n", dir)
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
		fmt.Printf("## cmd: %s\n", cmdfmt.Sprint(c.raw))
	}
	if len(c.raw) <= 0 {
		return nil, fmt.Errorf(cr.RedString("raw cmd is nil"))
	}
	if c.Execution != nil {
		return c.Execution(c.cmd, c.args...)
	}
	return nil, nil
}

func (c *CMD) Do() ([]byte, error) {
	if c.debug {
		fmt.Printf("## cmd: %s\n", cmdfmt.Sprint(c.raw))
	}
	if len(c.raw) <= 0 {
		return nil, fmt.Errorf(cr.RedString("raw cmd is nil"))
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
			goutils.CheckNoLogErr(err)
			if c.debug {
				fmt.Println(goutils.ToString(bs))
			}
			return bs, err
		case <-time.After(8e9):
			err = fmt.Errorf("Timeout")
		}
	}
	return nil, err
}

func (c *CMD) DoTimeout(dur time.Duration) ([]byte, error) {
	if c.debug {
		fmt.Printf("## cmd: %s\n", cmdfmt.Sprint(c.raw))
	}
	if len(c.raw) <= 0 {
		return nil, fmt.Errorf(cr.RedString("raw cmd is nil"))
	}
	var err error
	if c.Execution != nil {
		var bs []byte
		sh := make(chan bool)
		go func() {
			go func() {
				<-time.After(dur)
				runtime.Goexit()
			}()
			bs, err = c.Execution(c.cmd, c.args...)
			sh <- true
		}()
		select {
		case <-sh:
			goutils.CheckNoLogErr(err)
			if c.debug {
				fmt.Println(goutils.ToString(bs))
			}
			return bs, err
		case <-time.After(dur):
			err = fmt.Errorf("Timeout")
		}
	}
	return nil, err
}

func (c *CMD) Execute() *CMD {
	bs, err := c.DoNoTime()
	if err != nil {
		cr.Red(err.Error())
	}
	if c.debug {
		fmt.Print("## execute output: ")
		cr.Magenta("%s", goutils.ToString(bs))
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
	fmt.Println(goutils.ToString(bs))
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
