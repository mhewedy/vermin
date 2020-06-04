package db

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Box struct {
	CPU      string
	Mem      string
	DiskSize string
	MACAddr  string
}

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

func GetBoxInfo(vm string) (*Box, error) {
	var vb vbox
	b, _ := ioutil.ReadFile(GetVMPath(vm) + "/" + vm + ".vbox")
	err := xml.Unmarshal(b, &vb)

	if err != nil {
		return nil, err
	}

	cpuCount := vb.Machine.Hardware.CPU.Count
	if len(cpuCount) == 0 {
		cpuCount = "1"
	}
	return &Box{
		CPU:      cpuCount,
		Mem:      vb.Machine.Hardware.Memory.RAMSize,
		DiskSize: getDiskSizeInGB(vm, vb.Machine.MediaRegistry.HardDisks.HardDisk.Location),
		MACAddr:  vb.Machine.Hardware.Network.Adapter.MACAddress,
	}, nil
}

func getDiskSizeInGB(vm string, hdLocation string) string {
	stat, err := os.Stat(GetVMPath(vm) + string(os.PathSeparator) + hdLocation)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0))
}
