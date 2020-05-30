package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var (
	format = "%-15s%-27s%-10s%-10s%-13s%s\n"
	header = fmt.Sprintf(format, "VM NAME", "IMAGE", "CPUS", "MEM", "DISK", "TAGS")
)

// --- vmInfo types

type vmInfo struct {
	name  string
	image string
	box   db.Box
	disk  string
	tags  string
}

func (v vmInfo) String() string {
	return fmt.Sprintf(format, v.name, v.image, v.box.CPU, v.box.Mem, v.disk, v.tags)
}

type vmInfoList []vmInfo

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

// --- filter types

type filter struct {
	name  string
	value string
}

func (f filter) apply(list vmInfoList) vmInfoList {
	filtered := make(vmInfoList, 0)
	// fields should match name of fields from struct vmInfo
	fields := []string{"name", "image", "tags"}

	for _, e := range list {
		rv := reflect.ValueOf(e)

		for _, field := range fields {
			if f.name == field && strings.Contains(rv.FieldByName(field).String(), f.value) {
				filtered = append(filtered, e)
			}
		}
	}
	return filtered
}

type filters []filter

func (f filters) apply(list vmInfoList) vmInfoList {
	if len(f) == 0 {
		return list
	}
	for _, filter := range f {
		list = filter.apply(list)
	}
	return list
}

// --- package functions

func Ps(all bool, f []string) (string, error) {
	filters, err := parseFilters(f)
	if err != nil {
		return "", err
	}

	vms, err := List(all)
	if err != nil {
		return "", err
	}
	return getVMInfoList(vms, filters), nil
}

func parseFilters(filters []string) ([]filter, error) {
	if len(filters) == 0 {
		return nil, nil
	}

	var out = make([]filter, len(filters))

	for i, f := range filters {
		parts := strings.Split(f, "=")
		if len(parts) != 2 || len(parts[1]) == 0 {
			return nil, fmt.Errorf("Failed to parse fitler: %s\n", f)
		}
		out[i] = filter{name: parts[0], value: parts[1]}
	}
	return out, nil
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

func getVMInfoList(vms []string, filters filters) string {

	if len(vms) == 0 {
		return header
	}

	infoList := make(vmInfoList, len(vms))
	var wg sync.WaitGroup
	wg.Add(len(vms))

	for i, vmName := range vms {
		go func(vm string, i int) {
			infoList[i] = getVMInfo(vm)
			wg.Done()
		}(vmName, i)
	}
	wg.Wait()

	infoList = filters.apply(infoList)

	return infoList.String()
}

func getVMInfo(vm string) vmInfo {
	if _, err := os.Stat(db.GetVMPath(vm)); os.IsNotExist(err) {
		return vmInfo{}
	}

	box, _ := db.GetBoxInfo(vm)
	disk := getDiskSizeInGB(vm, box.HDLocation)
	vmdb, _ := db.Load(vm)

	return vmInfo{
		name:  vm,
		image: vmdb.Image,
		box:   *box,
		disk:  disk,
		tags:  strings.Join(vmdb.Tags, ", "),
	}
}

func getDiskSizeInGB(vm string, hdLocation string) string {
	stat, err := os.Stat(db.GetVMPath(vm) + string(os.PathSeparator) + hdLocation)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0))
}
