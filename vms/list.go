package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/db/info"
	"os"
	"sort"
	"strings"
)

var (
	format = "%-15s%-27s%-10s%-10s%-13s%s\n"
	header = fmt.Sprintf(format, "VM NAME", "IMAGE", "CPUS", "MEM", "DISK", "TAGS")
)

type vmInfo struct {
	name  string
	image string
	box   *info.Box
	disk  string
	tags  string
}

func (v *vmInfo) String() string {
	return fmt.Sprintf(format, v.name, v.image, v.box.CPU, v.box.Mem, v.disk, v.tags)
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

	box, _ := info.GetBoxInfo(vm)
	disk := getDiskSizeInGB(vm, box.HDLocation)
	image, _ := db.ReadImageData(vm)
	tags, _ := db.ReadTags(vm)

	return &vmInfo{
		name:  vm,
		image: image,
		box:   box,
		disk:  disk,
		tags:  tags,
	}
}

func getDiskSizeInGB(vm string, hdLocation string) string {
	stat, err := os.Stat(db.GetVMPath(vm) + string(os.PathSeparator) + hdLocation)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0))
}
