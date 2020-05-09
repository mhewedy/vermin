
[![Build Status](https://github.com/mhewedy/vermin/workflows/Go/badge.svg)](https://github.com/mhewedy/vermin/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhewedy/vermin)](https://goreportcard.com/report/github.com/mhewedy/vermin)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# vermin

### The Smart Virtual Machines manager

Create, control and connect to VirtualBox VM instances.

<p>
	<img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/logo.png"  alt="logo" width="30%"/>
</p>
      
      
<a href="https://asciinema.org/a/327940?speed=2&autoplay=1&cols=150&rows=35&size=medium"><img src="https://asciinema.org/a/327940.png" width="836"/></a>

### Menu

- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
	- [Automatic Installation](#Automatic-installation)
	- [Manual installation](#Manual-installation)
   -  [Build from source](#Build-from-source)
- [Use cases](#Use-cases)
- [Usage](#Usage)
	- [Create a new VM](#Create-a-new-VM)
	- [List VMs](#List-VMs)
	- [Start VM](#Start-VM)
	- [SSH into VM](#SSH-into-VM)
	- [Stop VM](#Stop-VM)
	- [Remove VM](#Remove-VM)
	- [Transfer Files](#Transfer-Files)
	- [Port Forward](#Port-Forward)
- [Why not Vagrant](#Why-not-Vagrant)
- [Contributors](#Contributors)
----
## Prerequisites
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads)

## Installation
#### Automatic installation:
To install or update for macos and linux:
```shell script
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mhewedy/vermin/master/install.sh)"
```
To install or update on windows (PowerShell):
```
# Should run as Adminstarator
C:\> iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/mhewedy/vermin/master/install.ps1'))
```

#### Manual installation:

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

#### Build from Source:
Download the latest released source code archive file from [releases](https://github.com/mhewedy/vermin/releases/latest) then unzip:
```bash
go build
```
You can build using golang docker image:
```bash
# replace window by linux or darwin depending on your OS
docker run -it -v $(pwd):/go -e GOPATH='' -e GOOS='windows' golang:latest go build
``` 
## Use cases:
Vermin can be used when you need an easy way to obtain a Linux up and running in minutes.

For example, if you want to have an environment to try .NET Core and you don't want to mess with your own WSL installation, so you can create a VM to do whatever you want to do then remove it.

Or if you want to try to install a Kafka cluster, and you need something more than just a docker container, so you can work with its different commands or want to practice setting up a cluster manually.

Also, you can check [Why not Vagrant](#Why-not-Vagrant) section.

## Usage:
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
  mount       Mount host path inside the VM
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

#### Create a new VM
Use the following command to create a VM

```shell script
$ vermin create <image name>
# example
$ vermin create ubuntu/focal
```
Or in case you want to create and provision the VM: (see [sample_init_bionic.sh](https://github.com/mhewedy/vermin/blob/master/etc/samples-provision/sample_init_bionic.sh) for sample provision script)
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

> To get the most updated image list (along with images locations) use the -p flag `vermin images -p`. this will not affect cached images. it will only get the most updated image list (image names along with thier remote locations).

#### List VMs
```shell script
$ vermin ps
VM NAME        IMAGE                      CPUS      MEM       DISK         TAGS
vm_01          alpine/3.11                1         1024      0.8GB
vm_02          ubuntu/focal               1         1024      2.6GB
vm_03          centos/8                   1         1024      2.0GB
```

#### Start VM
```shell script
$ vermin start vm_01
```

#### SSH into VM
```shell script
$ vermin ssh vm_03
```

#### Stop VM
```shell script
$ vermin stop vm_03
```

#### Remove VM
Will stop and remove listed VMs
```shell script
$ vermin rm vm_03
```

#### Transfer Files:
You can transfer files between host machine and VM.

To copy remote file on VM to you local host in the current path:
```shell script
$ vermin cp vm_01 --remote-file /path/to/file/on/vm
```

To copy local file from your host to the VM's home directory:
```shell script
$ vermin cp vm_01 --local-file /path/to/file/on/host
```

#### Port Forward:
forward ports from VM to local host (all ports from 8080 to 8090):
```shell script
$ vermin port vm_01 8080-8090
```

## Why not Vagrant:
* **Vagrant** uses a `Vagrantfile` which I think is most suited to be source-controlled inside `git`  , and for some use case it is an overhead to create and maintain such file. In such cases **Vermin** come to the rescue. 
* **Vermin** is a single binary file that can be easily installed and removed.


## Contributors
Special thanks to @aseldesouky for contributing the logo
