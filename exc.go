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
	originPath string
	raw        string
	cmd        string
	args       []string
	debug      bool
	Execution  func(cmd string, args ...string) ([]byte, error)
}

var (
	cmdfmt = cr.New(cr.FgHiGreen, cr.Bold, cr.BgHiBlack)
)

func NewCMDf(format string, args ...interface{}) *CMD {
	return NewCMD(fmt.Sprintf(format, args...))
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

	newCMD.originPath = originPath()
	return newCMD
}

func Bash(cmd string) *CMD {
	return &CMD{raw: cmd, cmd: cmd, Execution: BashExecution}
}

func (c *CMD) Debug(debug ...bool) *CMD {
	if len(debug) > 0 {
		c.debug = debug[0]
		return c
	}
	c.debug = !c.debug
	return c
}

func (c *CMD) deferFunc() {
	if c.originPath != "" {
		c.Cd(c.originPath)
	}
}

func (c *CMD) Set(cmd string) {
	c.raw, c.cmd, c.Execution = cmd, cmd, DefaultExecution
}

func originPath() string {
	dir, err := os.Getwd()
	goutils.LogCheckErr(err)
	return dir
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
	cr.Cyan("= = = [ %s ] = = =\n", originPath())
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
	defer c.deferFunc()
	if c.debug {
		cmdfmt.Print(c.raw)
		fmt.Println()
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
	defer c.deferFunc()
	if c.debug {
		cmdfmt.Print(c.raw)
		fmt.Println()
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
	defer c.deferFunc()
	if c.debug {
		cmdfmt.Print(c.raw)
		fmt.Println()
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
	if c.debug && len(bs) > 0 {
		cr.Magenta("%s", goutils.ToString(bs))
	}
	return c
}

func (c *CMD) Exec(debug bool) *CMD {
	bs, err := c.DoNoTime()
	if err != nil {
		cr.Red(err.Error())
	}
	if debug && len(bs) > 0 {
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
