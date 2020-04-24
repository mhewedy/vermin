package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// ExecuteP execute and show progress
func ExecuteP(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}

	go func() {
		for {
			fmt.Print(".")
			time.Sleep(3 * time.Second)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return "", errors.New(string(stderr.Bytes()))
	}
	fmt.Println()

	return string(stdout.Bytes()), nil
}

// Execute execute commands and return output as string or error
func Execute(command string, args ...string) (string, error) {
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

// Run runs command and return nothing
func Run(command string, args ...string) {
	_ = exec.Command(command, args...).Run()
}
