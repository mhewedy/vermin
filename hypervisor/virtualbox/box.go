package virtualbox

import (
	"encoding/xml"
	"fmt"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor/base"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
			Network struct {
				Adapter struct {
					MACAddress string `xml:"MACAddress,attr"`
				} `xml:"Adapter"`
			} `xml:"Network"`
		} `xml:"Hardware"`
	} `xml:"Machine"`
}

func getBoxInfo(vm string) (*base.Box, error) {
	var vb vbox
	b, _ := ioutil.ReadFile(db.GetVMPath(vm) + "/" + vm + ".vbox")
	err := xml.Unmarshal(b, &vb)

	if err != nil {
		return nil, err
	}

	cpuCount := vb.Machine.Hardware.CPU.Count
	if len(cpuCount) == 0 {
		cpuCount = "1"
	}
	return &base.Box{
		CPU:      cpuCount,
		Mem:      vb.Machine.Hardware.Memory.RAMSize,
		DiskSize: getDiskSizeInGB(vm, vb.Machine.MediaRegistry.HardDisks.HardDisk.Location),
		MACAddr:  vb.Machine.Hardware.Network.Adapter.MACAddress,
	}, nil
}

func getDiskSizeInGB(vm string, hdLocation string) string {
	stat, err := os.Stat(filepath.Join(db.GetVMPath(vm), hdLocation))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0))
}
