package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/hypervisor/base"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var (
	format = "%-15s%-30s%-10s%-10s%-13s%s\n"
	header = fmt.Sprintf(format, "VM NAME", "IMAGE", "CPUS", "MEM", "DISK", "TAGS")
)

// --- vmInfo types

type vmInfo struct {
	name  string
	image string
	box   base.Box
	tags  string
}

func (v vmInfo) String() string {
	return fmt.Sprintf(format, v.name, v.image, v.box.CPU, v.box.Mem, v.box.Disk.Size, v.tags)
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
			if f.name == field && strings.Contains(strings.ToLower(rv.FieldByName(field).String()), f.value) {
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

	vms, err := hypervisor.List(all)
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
		out[i] = filter{name: strings.ToLower(parts[0]), value: strings.ToLower(parts[1])}
	}
	return out, nil
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

	box, _ := hypervisor.GetBoxInfo(vm)
	vmdb, _ := db.Load(vm)

	return vmInfo{
		name:  vm,
		image: strings.Replace(vmdb.Image, "vagrant/", "", 1),
		box:   *box,
		tags:  strings.Join(vmdb.Tags, ", "),
	}
}
