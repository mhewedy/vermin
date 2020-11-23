package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/progress"
	"os"
	"os/exec"
	"runtime"
)

type Cmd struct {
	Command string
	Args    []string
}

func (c *Cmd) Call() (string, error) {
	return c.call(false, "")
}

func (c *Cmd) CallWithProgress(msg string) (string, error) {
	return c.call(true, msg)
}

func (c *Cmd) call(showProgress bool, msg string) (string, error) {

	if runtime.GOOS == "windows" {
		c.Args = prepend(c.Args, c.Command)
		c.Command = "powershell"
	}

	c.log()

	cmd := exec.Command(c.Command, c.Args...)

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
func (c *Cmd) Interact() error {
	if runtime.GOOS == "windows" {
		c.Args = prepend(c.Args, c.Command)
		c.Command = "powershell"
	}

	c.log()

	cmd := exec.Command(c.Command, c.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// Run runs Command and return nothing
func (c *Cmd) Run() error {
	if runtime.GOOS == "windows" {
		c.Args = prepend(c.Args, c.Command)
		c.Command = "powershell"
	}

	c.log()

	return exec.Command(c.Command, c.Args...).Run()
}

func prepend(x []string, y string) []string {
	return append([]string{y}, x...)
}

func (c *Cmd) log() {
	if _, ok := os.LookupEnv("VERMIN_DEBUG_CMD"); ok {
		fmt.Print("$ ", c.Command, " ")
		for _, arg := range c.Args {
			fmt.Print(arg, " ")
		}
		fmt.Println()
	}
}
