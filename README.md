<p align="center">
  <img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/logo.png"  alt="logo" width="70%"/> </center>
</p>

[![Build Status](https://github.com/mhewedy/vermin/workflows/Go/badge.svg)](https://github.com/mhewedy/vermin/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhewedy/vermin)](https://goreportcard.com/report/github.com/mhewedy/vermin)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


# vermin
### The Smart Virtual Machines manager    

Create, control and connect to VirtualBox VM instances.

----
## Prerequisites
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads)

## Installation
For macos and linux:
```shell script
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mhewedy/vermin/master/install.sh)"
```
For windows (PowerShell):
```
# Should run as Adminstarator
C:\> iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/mhewedy/vermin/master/install.ps1'))
```

#### Manual installation:

> It is recommended to use the method above to install vermin, However If you prefer to do manual installation then you need to follow these steps:

1. Download the binary matching your OS from [releases](https://github.com/mhewedy/vermin/releases/latest)
2. create the following directory structure in your home dir:
```
$HOME/.vermin
         ├── images
         └── vms
```
3. Download [vermin private key](https://raw.githubusercontent.com/mhewedy/vermin/master/etc/keys/vermin_rsa) into `$HOME/.vermin/vermin_rsa`
4. On windows, you need to add `C:\Program Files\Oracle\VirtualBox` into you PATH environment variable.

#### Build from source:
Download the latest released source code archieve file from [releases](https://github.com/mhewedy/vermin/releases/latest) then:
```
cd vermin && go build main.go
```

## Usage:

<p align="center">
  <img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/vermin-v0.35-demo.gif"  alt="demo" width="120%"/> </center>
</p>


```text
$ vermin
Create, control and connect to VirtualBox VM instances

You can start using vermin by creating a vm from an image.
To list all images available:
$ vermin images

Then you can create a vm using:
$ vermin create <image>

Usage:
  vermin [command]

Available Commands:
  completion  Generates completion scripts (Bash, Zsh and PowerShell)
  cp          Copy files between host and VM
  create      Create VM from an image
  help        Help about any command
  images      List all available images
  ip          Show IP address for a running VM
  port        Forward port(s) from a VM to host
  ps          List VMs
  rm          Remove a VM
  ssh         ssh into a running VM
  start       Start a VM
  stop        Stop a VM
  tag         Tag a VM

Flags:
  -h, --help      help for vermin
  -v, --version   version for vermin

Use "vermin [command] --help" for more information about a command.
```

### Create a new VM
Use the following command to create a VM

```shell script
$ vermin create <image name>
# example
$ vermin create ubuntu/focal
```
Or in case you want to create and provision the VM: (see [sample.sh](https://github.com/mhewedy/vermin/blob/master/etc/samples-provision/sample.sh) for sample provision script)
```shell script
$ vermin create <image name> /path/to/provison.sh 
# example
$ vermin create ubuntu/focal ~/sample.sh -cpus 1 -mem 512
```

To get list of all available images use:
```shell script
$ vermin images
ubuntu/focal	(cached)
centos/8
```
> The *cached* flag means, the image has been already downloaded and cached before.

### List all running VMs
```shell script
$ vermin ps
VM NAME		IMAGE				CPUS	MEM	TAGS
vm_01		ubuntu/focal			1	1024
```

### Start one or more VM
```shell script
$ vermin start vm_01
```

### ssh into a VM
```shell script
$ vermin ssh vm_03
```

### Stop one or more VMs
```shell script
$ vermin stop vm_03
```

### Remove one or more VMs
Will stop and remove listed VMs
```shell script
$ vermin rm vm_03
```

### Copy files:
Copy remote file on VM to you local host in the current path:
```shell script
$ vermin cp vm_01 --remote-file /path/to/file/on/vm
```

Copy local file from your host to the VM's home directory:
```shell script
$ vermin cp vm_01 --local-file /path/to/file/on/host
```

### Port forward:
forward ports from VM to local host (all ports from 8080 to 8090):
```shell script
$ vermin port vm_01 8080-8090
```

## Why not Vagrant:
* **Vagrant** uses a `Vagrantfile` which I think is most suited to be source-controlled, and for my case it is an overhead to maintain such file for each vm I want to create. (like create k8s cluster, etc...), I want kind of global accessibility.
