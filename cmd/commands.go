package cmd

import (
	"github.com/mhewedy/vermin/db"
	"runtime"
)

func Scp(extraArgs ...string) *Cmd {
	args := []string{"-q", "-r",
		"-i", db.GetPrivateKeyPath(),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}

	return &Cmd{
		Command: "scp",
		Args:    append(args, extraArgs...),
	}
}

func Arp(args ...string) *Cmd {
	return &Cmd{
		Command: "arp",
		Args:    args,
	}
}

func Ssh(ipAddr string, extraArgs ...string) *Cmd {
	args := []string{"-i", db.GetPrivateKeyPath(),
		"-o", "StrictHostKeyChecking=no",
		"-o", "GlobalKnownHostsFile=/dev/null",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=error",
		db.GetUsername() + "@" + ipAddr}

	return &Cmd{
		Command: "ssh",
		Args:    append(args, extraArgs...),
	}
}

func Ping(ip string) *Cmd {
	if runtime.GOOS == "windows" {
		return &Cmd{
			Command: "ping",
			Args:    []string{"-n", "1", "-w", "0.1", ip},
		}
	} else {
		return &Cmd{
			Command: "ping",
			Args:    []string{"-c", "1", "-W", "0.1", ip},
		}
	}
}

func AnsiblePlaybook(ip string, playbook string) *Cmd {
	return &Cmd{
		Command: "ansible-playbook",
		Args: []string{
			"-i", ip + ",",
			"-e", "ansible_user=" + db.GetUsername(),
			"-e", "ansible_private_key_file=" + db.GetPrivateKeyPath(),
			"--ssh-common-args", "-o StrictHostKeyChecking=no -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=/dev/null",
			playbook,
		},
	}
}
