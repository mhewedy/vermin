package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	format = "%-15s%-27s%-10s%-10s%-13s%s\n"
	header = fmt.Sprintf(format, "VM NAME", "IMAGE", "CPUS", "MEM", "DISK", "TAGS")
)

type vmInfo struct {
	name  string
	image string
	box   *db.Box
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

	numVms := len(vms)
	if numVms == 0 {
		return header
	}

	infoList := make(vmInfoList, numVms)
	var wg sync.WaitGroup
	wg.Add(numVms)

	for i, vmName := range vms {
		go func(vm string, i int) {
			infoList[i] = getVMInfo(vm)
			wg.Done()
		}(vmName, i)
	}
	wg.Wait()

	return infoList.String()
}

func getVMInfo(vm string) *vmInfo {
	if _, err := os.Stat(db.GetVMPath(vm)); os.IsNotExist(err) {
		return nil
	}

	box, _ := db.GetBoxInfo(vm)
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
