package command

import (
	"fmt"
	"runtime"
)

func VBoxManage(args ...string) *Cmd {
	return &Cmd{
		command: "vboxmanage",
		args:    args,
	}
}

func Scp(args ...string) *Cmd {
	return &Cmd{
		command: "scp",
		args:    args,
	}
}

func Arp(args ...string) *Cmd {
	return &Cmd{
		command: "arp",
		args:    args,
	}
}

func Wget(url string, file string) *Cmd {
	if runtime.GOOS == "windows" {
		return &Cmd{
			command: fmt.Sprintf("(New-Object System.Net.WebClient).DownloadFile('%s', '%s')", url, file),
		}
	} else {
		return &Cmd{
			command: "wget",
			args: []string{
				"-O",
				file,
				url,
			},
		}
	}
}

func Ssh(args ...string) *Cmd {
	return &Cmd{
		command: "ssh",
		args:    args,
	}
}

func Ping(ip string) *Cmd {
	if runtime.GOOS == "windows" {
		return &Cmd{
			command: "ping",
			args:    []string{"-n", "1", "-w", "0.1", ip},
		}
	} else {
		return &Cmd{
			command: "ping",
			args:    []string{"-c", "1", "-W", "0.1", ip},
		}
	}
}
