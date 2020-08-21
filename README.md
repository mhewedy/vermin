
[![Build Status](https://github.com/mhewedy/vermin/workflows/Go/badge.svg)](https://github.com/mhewedy/vermin/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhewedy/vermin)](https://goreportcard.com/report/github.com/mhewedy/vermin)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


<img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/logo.png"  alt="logo" width="20%"/>

## The smart virtual machines manager (star the project if you find it useful)

<img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/vermin2x.gif" width="100%"/>




Table of Contents:

- [What is Vermin](#what-is-vermin)
- [Install Vermin](#install-vermin)
- [Usage](#Usage)
- [Contributors](#Contributors)
- [Why not Vagrant](#Why-not-Vagrant)
- [TODO](#TODO)
----
# What is Vermin
Vermin is a smart, simple and powerful command line tool for Linux, Windows and macOS. It's designed for developers/tester and others working in IT who want a fresh VM environment with a single command. It uses VirtualBox to run the VM. Vermin will fetch images on your behalf.

Vermin can be used when you need an easy way to obtain a Linux environment up and running in minutes.
For example:
* If you want to have an environment to try .NET Core and you don't want to mess with your local own WSL installation.
* If you want to try to install a Kafka cluster, and you need something more than just a docker container.

Vermin in Action:
* [Install CockroachDB cluster in a Virtual Machine](https://medium.com/swlh/install-cockroachdb-on-a-virtual-machine-2f25878fd70?source=friends_link&sk=52b4c1c16794f8d15943c8c48a7103b5)
* [Install Redis inside Ubuntu VM](https://medium.com/swlh/install-redis-inside-a-ubuntu-vm-d5022d42d8cc?source=friends_link&sk=b7073861f8050c5318683d1ebfcd800a)
* [Install Kubernetes cluster in Virtual Machines the easy way](https://medium.com/@mhewedy_46874/install-kubernetes-cluster-in-virtual-machines-the-easy-way-337ef0c4e37f?source=friends_link&sk=dbb40739c54c864d1bd2f779032b2de2)
* [Install Desktop environment on Ubuntu Server](https://medium.com/@mhewedy_46874/ubuntu-20-04-desktop-vm-using-vermin-764d20f43c4d?source=friends_link&sk=d78dd1b863aaa0ea8bd05dc3e681c7ec)

Also, you can check [Why not Vagrant](#Why-not-Vagrant) section.

# Install Vermin

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

You can start using vermin after installation using:

```shell script
$ vermin create <image name> | vagrant/<vagrant image>
# example using vagrant image
$ vermin create vagrant/hashicorp/bionic64
```

### For more info on the usage options see [vermin documentations website](https://mhewedy.github.io/vermin/).


# Why not Vagrant:
* **Vagrant** uses a `Vagrantfile` which I think is most suited to be source-controlled inside `git`  , and for some use case it is an overhead to create and maintain such file. In such cases **Vermin** come to the rescue. 
* **Vermin** is a single binary file that can be easily installed and upgraded.
* It is important to note that, starting from version `v0.94.0` **Vermin** can smoothly uses Vagrant Cloud images.
* Myself, I look (and try to achieve) to Vermin as a modern CLI (like docker/podman) for Vagrant Boxes.

# Contributors
Special thanks to [Ahmed Samir](https://github.com/aseldesouky) for contributing the logo.

# TODO
See [TODO.md](https://github.com/mhewedy/vermin/blob/master/TODO.md)
