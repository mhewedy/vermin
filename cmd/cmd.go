package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
)

// ExecuteP execute and show progress
func ExecuteP(command string, args ...string) (string, error) {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}

	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	q := PrintProgress()
	defer func() { *q <- true }()

	if err := cmd.Wait(); err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	return string(stdout.Bytes()), nil
}

// Execute execute commands and return output as string or error
func Execute(command string, args ...string) (string, error) {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}

	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	return string(stdout.Bytes()), nil
}

// Execute execute commands in interactive mode
func ExecuteI(command string, args ...string) error {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}

	cmd := exec.Command(command, args...)

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
func Run(command string, args ...string) error {
	if runtime.GOOS == "windows" {
		args = prepend(args, command)
		command = "powershell"
	}
	return exec.Command(command, args...).Run()
}

func prepend(x []string, y string) []string {
	x = append(x, "")
	copy(x[1:], x)
	x[0] = y
	return x
}
