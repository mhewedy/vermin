package command

import (
	"github.com/mhewedy/vermin/db"
	"runtime"
)

func VBoxManage(args ...string) *cmd {
	return &cmd{
		command: "vboxmanage",
		args:    args,
	}
}

func Scp(vmName string, extraArgs ...string) *cmd {
	args := []string{"-q", "-r",
		"-i", db.GetPrivateKeyPath(vmName),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}

	return &cmd{
		command: "scp",
		args:    append(args, extraArgs...),
	}
}

func Arp(args ...string) *cmd {
	return &cmd{
		command: "arp",
		args:    args,
	}
}

func Ssh(vmName, ipAddr string, extraArgs ...string) *cmd {
	args := []string{"-i", db.GetPrivateKeyPath(vmName),
		"-o", "StrictHostKeyChecking=no",
		"-o", "GlobalKnownHostsFile=/dev/null",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=error",
		db.GetUsername(vmName) + "@" + ipAddr}

	return &cmd{
		command: "ssh",
		args:    append(args, extraArgs...),
	}
}

func Ping(ip string) *cmd {
	if runtime.GOOS == "windows" {
		return &cmd{
			command: "ping",
			args:    []string{"-n", "1", "-w", "0.1", ip},
		}
	} else {
		return &cmd{
			command: "ping",
			args:    []string{"-c", "1", "-W", "0.1", ip},
		}
	}
}

func AnsiblePlaybook(vmName, ip string, playbook string) *cmd {
	return &cmd{
		command: "ansible-playbook",
		args: []string{
			"-i", ip + ",",
			"-e", "ansible_user=" + db.GetUsername(vmName),
			"-e", "ansible_private_key_file=" + db.GetPrivateKeyPath(vmName),
			"--ssh-common-args", "-o StrictHostKeyChecking=no -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=/dev/null",
			playbook,
		},
	}
}

// Tar flags are (cxzf)
// Temporary, use go code instead
func Tar(args ...string) *cmd {
	return &cmd{
		command: "tar",
		args:    args,
	}
}
