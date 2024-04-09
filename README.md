
[![Build Status](https://github.com/mhewedy/vermin/workflows/Go/badge.svg)](https://github.com/mhewedy/vermin/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhewedy/vermin)](https://goreportcard.com/report/github.com/mhewedy/vermin)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<img src="https://raw.githubusercontent.com/mhewedy/vermin/master/etc/logo.png"  alt="logo" width="20%"/>

## The smart virtual machines manager
Table of Contents:

- [What is Vermin](#what-is-vermin)
- [Install Vermin](#install-vermin)
- [Usage](#Usage)
- [Contributors](#Contributors)
- [TODO](#TODO)

----

# What is Vermin
Vermin is a smart, simple and powerful command line tool for Linux, Windows and macOS. It's designed for developers/testers and others working in IT who want a fresh VM environment with a single command. It uses VirtualBox to run the VM. Vermin will fetch images on your behalf.

You can look to Vermin as a modern CLI for Vagrant Boxes.

Vermin can be used when you need an easy way to obtain a Linux environment up and running in minutes.

# Install Vermin

Vermin uses [VirtualBox v6.0 or later](https://www.virtualbox.org/wiki/Downloads) as the underlying hypervisor to create and run Virtual Machines. So you need to download and install it first.

To install/update on **macOS** and **Linux** run:

```shell script
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mhewedy/vermin/master/install.sh)"
```

To install/update on **Windows** (PowerShell) run:

```
# Should run as Administrator
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
  hypervisor  print the name of the detected hypervisor
  images      List remote and cached images
  ip          Show IP address for a running VM
  mount       Mount local filesystem inside the VM
  port        Forward port(s) from a VM to host
  ps          List VMs
  restart     Restart one or more VMs
  rm          Remove one or more VM
  rmi         Remove one or more Image
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

You can start using Vermin after installation using:

```shell script
$ vermin create <vagrant image name>

# example using ubuntu focal image
$ vermin create hashicorp/focal64

# also you can use rhel8 using:
$ vermin create generic/rhel8
```
You can use all [vagrant images](https://app.vagrantup.com/boxes/search).

_Vermin collects very simple usage data anonymously._

## Collaborators

<!-- readme: collaborators -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/mhewedy">
            <img src="https://avatars.githubusercontent.com/u/1086049?v=4" width="100;" alt="mhewedy"/>
            <br />
            <sub><b>Mohammad Hewedy</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/akhiljns">
            <img src="https://avatars.githubusercontent.com/u/22194681?v=4" width="100;" alt="akhiljns"/>
            <br />
            <sub><b>Akhil</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: collaborators -end -->

## Contributors

<!-- readme: contributors -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/mhewedy">
            <img src="https://avatars.githubusercontent.com/u/1086049?v=4" width="100;" alt="mhewedy"/>
            <br />
            <sub><b>Mohammad Hewedy</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/akhiljns">
            <img src="https://avatars.githubusercontent.com/u/22194681?v=4" width="100;" alt="akhiljns"/>
            <br />
            <sub><b>Akhil</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/dawidd6">
            <img src="https://avatars.githubusercontent.com/u/9713907?v=4" width="100;" alt="dawidd6"/>
            <br />
            <sub><b>Dawid Dziurla</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/gruz0">
            <img src="https://avatars.githubusercontent.com/u/335095?v=4" width="100;" alt="gruz0"/>
            <br />
            <sub><b>Alexander Kadyrov</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/aldarisbm">
            <img src="https://avatars.githubusercontent.com/u/32185409?v=4" width="100;" alt="aldarisbm"/>
            <br />
            <sub><b>Jose Berrio</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: contributors -end -->

Special thanks to [Ahmed Samir](https://github.com/aseldesouky) for contributing the logo.

