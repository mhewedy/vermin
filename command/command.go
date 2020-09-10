package command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/progress"
	"os"
	"os/exec"
	"runtime"
)

type cmd struct {
	command string
	args    []string
}

func (c *cmd) Call() (string, error) {
	return c.call(false, "")
}

func (c *cmd) CallWithProgress(msg string) (string, error) {
	return c.call(true, msg)
}

func (c *cmd) call(showProgress bool, msg string) (string, error) {

	if runtime.GOOS == "windows" {
		c.args = prepend(c.args, c.command)
		c.command = "powershell"
	}

	c.log()

	cmd := exec.Command(c.command, c.args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", errors.New(err.Error() + " " + string(stderr.Bytes()))
	}

	if showProgress {
		stop := progress.Show(msg, false)
		defer stop()
	}

	if err := cmd.Wait(); err != nil {
		return "", errors.New(err.Error() + " " + string(stderr.Bytes()))
	}

	return string(stdout.Bytes()), nil
}

// Execute execute commands in interactive mode
func (c *cmd) Interact() error {
	if runtime.GOOS == "windows" {
		c.args = prepend(c.args, c.command)
		c.command = "powershell"
	}

	c.log()

	cmd := exec.Command(c.command, c.args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// Run runs command and return nothing
func (c *cmd) Run() error {
	if runtime.GOOS == "windows" {
		c.args = prepend(c.args, c.command)
		c.command = "powershell"
	}

	c.log()

	return exec.Command(c.command, c.args...).Run()
}

func prepend(x []string, y string) []string {
	return append([]string{y}, x...)
}

func (c *cmd) log() {
	if _, ok := os.LookupEnv("VERMIN_DEBUG"); ok {
		fmt.Print("$ ", c.command, " ")
		for _, arg := range c.args {
			fmt.Print(arg, " ")
		}
		fmt.Println()
	}
}
