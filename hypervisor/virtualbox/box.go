package virtualbox

import (
	"encoding/xml"
	"fmt"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor/base"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type vbox struct {
	XMLName xml.Name `xml:"VirtualBox"`
	Machine struct {
		MediaRegistry struct {
			HardDisks struct {
				HardDisk []struct {
					UUID     string `xml:"uuid,attr"`
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
		StorageControllers struct {
			StorageController []struct {
				AttachedDevice []struct {
					Port  string `xml:"port,attr"`
					Image struct {
						Uuid string `xml:"uuid,attr"`
					} `xml:"Image"`
				} `xml:"AttachedDevice"`
			} `xml:"StorageController"`
		} `xml:"StorageControllers"`
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

	diskUUID, diskLocation, err := findDiskUUIDAndLocation(vb)

	return &base.Box{
		CPU: cpuCount,
		Mem: vb.Machine.Hardware.Memory.RAMSize,
		Disk: &base.Disk{
			Size:     getDiskSizeInGB(vm, diskLocation),
			Location: diskLocation,
			UUID:     diskUUID,
		},
		MACAddr: vb.Machine.Hardware.Network.Adapter.MACAddress,
	}, nil
}

func findDiskUUIDAndLocation(vb vbox) (string, string, error) {

	var uuid, location string

outer:
	for _, sc := range vb.Machine.StorageControllers.StorageController {
		if len(sc.AttachedDevice) > 0 {

			for _, ad := range sc.AttachedDevice {
				if ad.Port == "0" {

					uuid = ad.Image.Uuid
					break outer
				}
			}
		}
	}

	for _, hd := range vb.Machine.MediaRegistry.HardDisks.HardDisk {
		if hd.UUID == uuid {
			location = hd.Location
			break
		}
	}

	// clean uuid
	uuid = strings.TrimFunc(uuid, func(r rune) bool {
		return r == '{' || r == '}'
	})

	return uuid, location, nil
}

func getDiskSizeInGB(vm string, hdLocation string) string {
	stat, err := os.Stat(filepath.Join(db.GetVMPath(vm), hdLocation))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0))
}
