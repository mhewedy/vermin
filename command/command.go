package command

import (
	"bytes"
	"errors"
	"github.com/mhewedy/vermin/progress"
	"os"
	"os/exec"
	"runtime"
)

type Cmd struct {
	command string
	args    []string
}

func (c *Cmd) Call() (string, error) {
	return c.call(false, "")
}

func (c *Cmd) CallWithProgress(msg string) (string, error) {
	return c.call(true, msg)
}

func (c *Cmd) call(showProgress bool, msg string) (string, error) {

	if runtime.GOOS == "windows" {
		c.args = prepend(c.args, c.command)
		c.command = "powershell"
	}

	cmd := exec.Command(c.command, c.args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	if showProgress {
		stop := progress.Show(msg)
		defer stop()
	}

	if err := cmd.Wait(); err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	return string(stdout.Bytes()), nil
}

// Execute execute commands in interactive mode
func (c *Cmd) Interact() error {
	if runtime.GOOS == "windows" {
		c.args = prepend(c.args, c.command)
		c.command = "powershell"
	}

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
func (c *Cmd) Run() error {
	if runtime.GOOS == "windows" {
		c.args = prepend(c.args, c.command)
		c.command = "powershell"
	}
	return exec.Command(c.command, c.args...).Run()
}

func prepend(x []string, y string) []string {
	return append([]string{y}, x...)
}
