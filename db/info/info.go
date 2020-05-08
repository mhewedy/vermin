package info

import (
	"encoding/xml"
	"github.com/mhewedy/vermin/db"
	"io/ioutil"
)

type VBox struct {
	CPU        string
	Mem        string
	HDLocation string
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
		} `xml:"Hardware"`
	} `xml:"Machine"`
}

func Get(vm string) (*VBox, error) {
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
	return &VBox{
		CPU:        cpuCount,
		Mem:        vb.Machine.Hardware.Memory.RAMSize,
		HDLocation: vb.Machine.MediaRegistry.HardDisks.HardDisk.Location,
	}, nil
}
