package vms

import (
	"encoding/xml"
	"fmt"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/db"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type vmInfo struct {
	name  string
	image string
	cpu   int
	mem   int
	tags  string
}

func Ps(all bool) (string, error) {
	vms, err := List(all)
	if err != nil {
		return "", err
	}
	return getVMsInfo(vms), nil
}

// List return all vms that start with db.NamePrefix
func List(all bool) ([]string, error) {
	var args = [2]string{"list"}
	if all {
		args[1] = "vms"
	} else {
		args[1] = "runningvms"
	}

	r, err := cmd.Execute("vboxmanage", args[:]...)
	if err != nil {
		return nil, err
	}

	var vms []string
	fields := strings.Fields(r)

	for i := range fields {
		if i%2 == 0 {
			vmName := strings.ReplaceAll(fields[i], `"`, "")
			if strings.HasPrefix(vmName, db.NamePrefix) {
				vms = append(vms, vmName)
			}
		}
	}

	return vms, nil
}

// get get info about vms
func getVMsInfo(vms []string) string {

	if len(vms) == 0 {
		return "VM NAME\t\tIMAGE\t\t\t\tCPUS\tMEM\tTAGS"
	}

	ch := make(chan *vmInfo, len(vms))

	for _, vmName := range vms {
		go func(vm string) {
			ch <- getVMInfo(vm)
		}(vmName)
	}

	// collect from channel
	var i int
	out := make([]*vmInfo, 0)

	for {
		select {
		case vmInfo := <-ch:
			if vmInfo != nil {
				out = append(out, vmInfo)
			}
			i++
		}
		if i == len(vms) {
			break
		}
	}

	return asString(out)
}

func asString(vmInfos []*vmInfo) string {
	var out string

	sort.Slice(vmInfos, func(i, j int) bool {
		return vmInfos[i].name < vmInfos[j].name
	})
	out += fmt.Sprintln("VM NAME\t\tIMAGE\t\t\t\tCPUS\tMEM\tTAGS")
	for _, e := range vmInfos {
		out += fmt.Sprintf("%s\t\t%s\t\t\t%d\t%d\t%s\n", e.name, e.image, e.cpu, e.mem, e.tags)
	}

	return out
}

func getVMInfo(vm string) *vmInfo {

	if _, err := os.Stat(db.GetVMPath(vm)); os.IsNotExist(err) {
		return nil
	}

	c, m := getVMCpuAndMem(vm)
	cpu, _ := strconv.Atoi(c)
	mem, _ := strconv.Atoi(m)

	image := readFromVMDB(vm, db.Image, "\t")
	tags := readFromVMDB(vm, db.Tags, "\t")

	return &vmInfo{
		name:  vm,
		image: image,
		cpu:   cpu,
		mem:   mem,
		tags:  tags,
	}
}

func readFromVMDB(vm string, dbFile string, defaultValue string) string {
	b, _ := ioutil.ReadFile(db.GetVMPath(vm) + "/" + dbFile)
	v := strings.ReplaceAll(string(b), "\n", " ")
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func getVMCpuAndMem(vm string) (string, string) {
	type vbox struct {
		XMLName xml.Name `xml:"VirtualBox"`
		Machine struct {
			Hardware struct {
				CPU struct {
					Count string `xml:"count,attr"`
				} `xml:"CPU"`
				Memory struct {
					RAMSize string `xml:"RAMSize,attr"`
				} `xml:"Memory"`
			} `xml:"Hardware"`
		} `xml:"Machine"`
	}

	var vb vbox
	b, _ := ioutil.ReadFile(db.GetVMPath(vm) + "/" + vm + ".vbox")
	err := xml.Unmarshal(b, &vb)

	if err != nil {
		return "", ""
	}

	cpuCount := vb.Machine.Hardware.CPU.Count
	if len(cpuCount) == 0 {
		cpuCount = "1"
	}
	return cpuCount,
		vb.Machine.Hardware.Memory.RAMSize
}
