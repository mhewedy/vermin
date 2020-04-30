<p align="center">
  <img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/logo.png"  alt="logo" width="70%"/> </center>
</p>

[![Build Status](https://github.com/mhewedy/vermin/workflows/Go/badge.svg)](https://github.com/mhewedy/vermin/actions?query=workflow%3AGo)

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
For windows: comming soon

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
  completion  Generates completion scripts (bash and zsh)
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
  -h, --help   help for vermin

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
