package info

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"vermin/db"
)

type vmInfo struct {
	name  string
	image string
	cpu   int
	mem   int
	tags  string
}

// List list info about vms
//  vm_01 vm_02 vm_03
func List(vms []string) string {

	if len(vms) == 0 {
		return "VM NAME\t\tIMAGE\t\t\tCPU\t\tMEM\t\tTAGS"
	}

	ch := make(chan *vmInfo, len(vms))

	for _, vmName := range vms {
		go func(vm string) {
			ch <- getVMInfo(vm)
		}(vmName)
	}

	// collect from channel
	out := make([]*vmInfo, 0)

	var i int
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

	return printInfo(out)
}

func printInfo(vmInfos []*vmInfo) string {
	var out string

	sort.Slice(vmInfos, func(i, j int) bool {
		return vmInfos[i].name < vmInfos[j].name
	})
	out += fmt.Sprintln("VM NAME\t\tIMAGE\t\t\tCPU\t\tMEM\t\tTAGS")
	for _, e := range vmInfos {
		out += fmt.Sprintf("%s\t\t%s\t\t%d\t\t%d\t\t%s\n", e.name, e.image, e.cpu, e.mem, e.tags)
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