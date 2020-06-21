
[![Build Status](https://github.com/mhewedy/vermin/workflows/Go/badge.svg)](https://github.com/mhewedy/vermin/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhewedy/vermin)](https://goreportcard.com/report/github.com/mhewedy/vermin)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


<img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/logo.png"  alt="logo" width="20%"/>

## The smart virtual machines manager (star the project if you find it useful)

<img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/vermin2x.gif"/>




Table of Contents:

- [What is Vermin](#what-is-vermin)
- [Install Vermin](#install-vermin)
- [Usage](#Usage)
	- [Create a new VM](#Create-a-new-VM)
	- [List VMs](#List-VMs)
	- [Start VM](#Start-VM)
	- [SSH into VM](#SSH-into-VM)
	- [Stop VM](#Stop-VM)
	- [Remove VM](#Remove-VM)
	- [Transfer Files](#Transfer-Files)
	- [Port Forward](#Port-Forward)
	- [Commit VM](#Commit-VM) (Create your own images)
	- [And more](#Usage)
- [More installation options](#more-installation-options)
   	- [Manual installation](#Manual-installation)
   	- [Build from source](#Build-from-source)
- [Contributors](#Contributors)
- [Why not Vagrant](#Why-not-Vagrant)
- [TODO](#TODO)
----
# What is Vermin
Vermin is a smart, simple and powerful command line tool for Linux, Windows and macOS. It's designed for developers who want a fresh VM environment with a single command. It uses VirtualBox to run the VM. Vermin will fetch images in your behave.

Vermin can be used when you need an easy way to obtain a Linux up and running in minutes.
For example:
* If you want to have an environment to try .NET Core and you don't want to mess with your local own WSL installation.
* Or if you want to try to install a Kafka cluster, and you need something more than just a docker container.

Vermin in Action:
* [Install CockroachDB cluster in a Virtual Machine](https://medium.com/swlh/install-cockroachdb-on-a-virtual-machine-2f25878fd70?source=friends_link&sk=52b4c1c16794f8d15943c8c48a7103b5)
* [Install Redis inside Ubuntu VM](https://medium.com/swlh/install-redis-inside-a-ubuntu-vm-d5022d42d8cc?source=friends_link&sk=b7073861f8050c5318683d1ebfcd800a)
* [Install Kubernetes cluster in Virtual Machines the easy way](https://medium.com/@mhewedy_46874/install-kubernetes-cluster-in-virtual-machines-the-easy-way-337ef0c4e37f?source=friends_link&sk=dbb40739c54c864d1bd2f779032b2de2)
* [Install Desktop environment on Ubuntu Server](https://medium.com/@mhewedy_46874/ubuntu-20-04-desktop-vm-using-vermin-764d20f43c4d?source=friends_link&sk=d78dd1b863aaa0ea8bd05dc3e681c7ec)

Also, you can check [Why not Vagrant](#Why-not-Vagrant) section.

# Install Vermin 🚀

Vermin uses [VirtualBox](https://www.virtualbox.org/wiki/Downloads) as the underlying hypervisor to create and run Virtual Machines. So you need to download and install it first.

To install/update on **macos** and **linux** run:
```shell script
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mhewedy/vermin/master/install.sh)"
```
To install/update on **windows** (PowerShell) run:
```
# Should run as Adminstarator
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/mhewedy/vermin/master/install.ps1'))
```

# Usage:
```text
Create, control and connect to VirtualBox VM instances

Usage:
  vermin [command]

Examples:

You can use vermin by creating a VM from an image.

To list all images available:
$ vermin images

Then you can create a vm using:
$ vermin create <image>


Available Commands:
  commit      Commit a VM into a new Image
  completion  Generates shell completion scripts
  cp          Copy files/folders between a VM and the local filesystem or between two VMs
  create      Create a new VM
  exec        Run a command in a running VM
  gui         open the GUI for the VM
  help        Help about any command
  images      List remote and cached images
  ip          Show IP address for a running VM
  mount       Mount local filesystem inside the VM
  port        Forward port(s) from a VM to host
  ps          List VMs
  restart     Restart one or more VMs
  rm          Remove one or more VM
  ssh         ssh into a running VM
  start       Start one or more stopped VMs
  stop        Stop one or more running VMs
  tag         Add or remove tag to a VM
  update      Update configuration of a VM

Flags:
  -h, --help      help for vermin
  -v, --version   version for vermin

Use "vermin [command] --help" for more information about a command.
```

## Create a new VM
Use the following command to create a VM

```shell script
$ vermin create <image name>
# example
$ vermin create ubuntu/focal
```
Or in case you want to create and provision the VM: (example [Sample provision scripts](https://github.com/mhewedy/vermin/blob/master/etc/samples-provision))
```shell script
$ vermin create <image name> /path/to/provison.sh 
# example
$ vermin create ubuntu/focal ~/sample.sh -cpus 1 -mem 512
```

To get list of all available images use:
```shell script
$ vermin images
IMAGE NAME               CACHED              DISK
centos/8                 true                0.8GB
ubuntu/focal             true                1.1GB
alpine/3.11
```
> The *cached* flag means, the image has been already downloaded and cached before.

> To get the most updated image list (along with images locations) use the -p flag `vermin images -p`. this will not affect cached images. it will only get the most updated image list (image names along with thier remote locations).

## List VMs
```shell script
$ vermin ps
VM NAME        IMAGE                      CPUS      MEM       DISK         TAGS
vm_01          alpine/3.11                1         1024      0.8GB
vm_02          ubuntu/focal               1         1024      2.6GB
vm_03          centos/8                   1         1024      2.0GB
```

## Start VM
```shell script
$ vermin start vm_01
```

## SSH into VM
```shell script
$ vermin ssh vm_03
```

## Stop VM
```shell script
$ vermin stop vm_03
```

## Remove VM
Will stop and remove listed VMs
```shell script
$ vermin rm vm_03
```

## Transfer Files:
You can transfer files between host machine and VM.

To copy a remote file on a VM to you local host in the current path:
```shell script
$ vermin cp vm_01:/path/to/file/on/vm .
```

To copy a local file from your host filesystem to the VM's home directory:
```shell script
$ vermin cp /path/to/file/on/host vm_01:~
```

Copy file.txt from vm_01 home dir to the vm_02 /tmp dir
```shell script
$ vermin cp vm_01:~/file.txt vm_02:/tmp
```

## Port Forward:
forward ports from VM to local host (all ports from 8080 to 8090):
```shell script
$ vermin port vm_01 8080-8090
```

## Commit VM:
You can commit a stopped VM into an image, so later you can create VMs from that image:
```shell script
$ vermin commit vm_01 elk/elastic
```
Then you can create as many VMs as you want from the committed image:
```shell script
$ vermin create elk/elastic
```

To remove an image at any time later use (removing an image doesn't remove VMs created from it):
```shell script
$ vermin image rm elk/elastic
```

# More installation options:
## Manual installation:

> It is recommended to use the [automatic method](#Automatic-installation) to install vermin, However If you prefer to do manual installation then you need to follow these steps:

1. Download the binary matching your OS from [releases](https://github.com/mhewedy/vermin/releases/latest) unzip it and preferably put it in your PATH 
2. create the following directory structure in your home dir:
```
$HOME/.vermin
         ├── images
         └── vms
```
3. Download [vermin private key](https://raw.githubusercontent.com/mhewedy/vermin/master/etc/keys/vermin_rsa) into `$HOME/.vermin/vermin_rsa`
4. On windows, you need to add `C:\Program Files\Oracle\VirtualBox` into you PATH.

## Build from Source:
Download the latest released source code archive file from [releases](https://github.com/mhewedy/vermin/releases/latest) then unzip:
```bash
go build
```
You can build using golang docker image:
```bash
# replace window by linux or darwin depending on your OS
docker run -it -v $(pwd):/go -e GOPATH='' -e GOOS='windows' golang:latest go build
``` 
## Arch Linux Install
```bash
yay -S vermin-bin
```
```bash
trizen -S vermin-bin
```
# Why not Vagrant:
* **Vagrant** uses a `Vagrantfile` which I think is most suited to be source-controlled inside `git`  , and for some use case it is an overhead to create and maintain such file. In such cases **Vermin** come to the rescue. 
* **Vermin** is a single binary file that can be easily installed and upgraded.

# Contributors
Special thanks to [Ahmed Samir](https://github.com/aseldesouky) for contributing the logo.

# TODO
See [TODO.md](https://github.com/mhewedy/vermin/blob/master/TODO.md)
