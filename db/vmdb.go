package db

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Decided to use gob package instead of protobuf after reading this https://blog.golang.org/gob

const vmdbFilename = "vmdb"

type VMDB struct {
	Image string
	Tags  []string
}

func Load(vm string) (VMDB, error) {
	migrate(vm)
	file, err := os.Open(GetVMPath(vm) + "/" + vmdbFilename)
	if err != nil {
		return VMDB{}, err
	}
	defer file.Close()

	return load(file)
}

func SetImage(vm string, image string) error {
	migrate(vm)
	return override(vm, func(vmdb *VMDB) {
		vmdb.Image = image
	})
}

func AddTag(vm string, tag string) error {
	migrate(vm)
	return override(vm, func(vmdb *VMDB) {
		vmdb.Tags = append(vmdb.Tags, tag)
	})
}

func RemoveTag(vm string, tag string) error {
	migrate(vm)
	return override(vm, func(vmdb *VMDB) {
		for i, v := range vmdb.Tags {
			if v == tag {
				vmdb.Tags = append(vmdb.Tags[:i], vmdb.Tags[i+1:]...)
				break
			}
		}
	})
}

// Override the VMDB for the vm passed as the first parameter
// by applying the fn function on the data loaded from disk
// and then write back to the disk
func override(vm string, fn func(vmdb *VMDB)) error {
	file, err := os.OpenFile(GetVMPath(vm)+"/"+vmdbFilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	vmdb, err := load(file)
	if err != nil {
		return err
	}

	fn(&vmdb)

	// clear the file before save
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	// serialize
	if err := gob.NewEncoder(file).Encode(vmdb); err != nil {
		return err
	}

	return nil
}

func load(file *os.File) (VMDB, error) {
	var vmdb VMDB
	if err := gob.NewDecoder(file).Decode(&vmdb); err != nil && io.EOF != err {
		return VMDB{}, err
	}
	return vmdb, nil
}

// --- migrate old db structure
func migrate(vm string) {
	if _, err := os.Stat(GetVMPath(vm) + "/" + vmdbFilename); err == nil {
		return // already migrated
	}
	_ = override(vm, func(vmdb *VMDB) {
		b, err := ioutil.ReadFile(GetVMPath(vm) + "/image")
		if err == nil {
			vmdb.Image = string(b)
		}

		b, err = ioutil.ReadFile(GetVMPath(vm) + "/tags")
		if err == nil {
			tags := string(b)
			split := strings.Split(tags, "\n")
			vmdb.Tags = split[:len(split)-1]
		}
	})
}
