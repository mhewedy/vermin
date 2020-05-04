package db

import (
	"io"
	"io/ioutil"
	"os"
)

func appendToFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func readFromFile(vm string, dbFile string) (string, error) {
	b, err := ioutil.ReadFile(GetVMPath(vm) + "/" + dbFile)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
