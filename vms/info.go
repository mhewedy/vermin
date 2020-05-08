package vms

import (
	"encoding/xml"
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	format = "%-15s%-27s%-10s%-10s%-13s%s\n"
	header = fmt.Sprintf(format, "VM NAME", "IMAGE", "CPUS", "MEM", "DISK", "TAGS")
)

type vmInfo struct {
	name    string
	image   string
	hwSpecs hwSpecs
	disk    string
	tags    string
}

type hwSpecs struct {
	cpu        string
	mem        string
	hdLocation string
}

func (v *vmInfo) String() string {
	return fmt.Sprintf(format, v.name, v.image, v.hwSpecs.cpu, v.hwSpecs.mem, v.disk, v.tags)
}

type vmInfoList []*vmInfo

func (l vmInfoList) String() string {
	var out string

	sort.Slice(l, func(i, j int) bool {
		return l[i].name < l[j].name
	})
	out += header
	for _, e := range l {
		out += e.String()
	}

	return out
}

func Ps(all bool) (string, error) {
	vms, err := List(all)
	if err != nil {
		return "", err
	}
	return getVMInfoList(vms), nil
}

// List return all vms that start with db.VMNamePrefix
func List(all bool) ([]string, error) {
	var args = [2]string{"list"}
	if all {
		args[1] = "vms"
	} else {
		args[1] = "runningvms"
	}

	r, err := command.VBoxManage(args[:]...).Call()
	if err != nil {
		return nil, err
	}

	var vms []string
	fields := strings.Fields(r)

	for i := range fields {
		if i%2 == 0 {
			vmName := strings.ReplaceAll(fields[i], `"`, "")
			if strings.HasPrefix(vmName, db.VMNamePrefix) {
				vms = append(vms, vmName)
			}
		}
	}

	return vms, nil
}

// get get info about vms
func getVMInfoList(vms []string) string {

	if len(vms) == 0 {
		return header
	}

	ch := make(chan *vmInfo, len(vms))

	for _, vmName := range vms {
		go func(vm string) {
			ch <- getVMInfo(vm)
		}(vmName)
	}

	// collect from channel
	var i int
	infoList := make(vmInfoList, 0)

	for {
		select {
		case vmInfo := <-ch:
			if vmInfo != nil {
				infoList = append(infoList, vmInfo)
			}
			i++
		}
		if i == len(vms) {
			break
		}
	}

	return infoList.String()
}

func getVMInfo(vm string) *vmInfo {

	if _, err := os.Stat(db.GetVMPath(vm)); os.IsNotExist(err) {
		return nil
	}

	hw := getHWSpecs(vm)
	disk := getDiskSizeGB(vm, hw.hdLocation)
	image, _ := db.ReadImageData(vm)
	tags, _ := db.ReadTags(vm)

	return &vmInfo{
		name:    vm,
		image:   image,
		hwSpecs: hw,
		disk:    disk,
		tags:    tags,
	}
}

func getDiskSizeGB(vm string, hdLocation string) string {
	stat, err := os.Stat(db.GetVMPath(vm) + string(os.PathSeparator) + hdLocation)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0))
}

func getHWSpecs(vm string) hwSpecs {
	type vbox struct {
		XMLName xml.Name `xml:"VirtualBox"`
		Machine struct {
			MediaRegistry struct {
				HardDisks struct {
					HardDisk struct {
						Location string `xml:"location,attr"`
					} `xml:"HardDisk"`
				} `xml:"HardDisks"`
			} `xml:"MediaRegistry"`
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
		return hwSpecs{}
	}

	cpuCount := vb.Machine.Hardware.CPU.Count
	if len(cpuCount) == 0 {
		cpuCount = "1"
	}
	return hwSpecs{
		cpu:        cpuCount,
		mem:        vb.Machine.Hardware.Memory.RAMSize,
		hdLocation: vb.Machine.MediaRegistry.HardDisks.HardDisk.Location,
	}
}
