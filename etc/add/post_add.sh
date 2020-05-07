#!/bin/bash

function install_pkg() {
  declare -A osInfo;
  osInfo[/etc/debian_version]="apt-get install -y"
  osInfo[/etc/alpine-release]="apk --update add"
  osInfo[/etc/centos-release]="yum install -y"
  osInfo[/etc/fedora-release]="dnf install -y"

  # shellcheck disable=SC2068
  for f in ${!osInfo[@]}
  do
      if [[ -f $f ]];then
          package_manager=${osInfo[$f]}
      fi
  done

  sudo ${package_manager} "$1"
}

function check_user() {
  if ! id "vermin" >/dev/null 2>&1; then
    echo "vermin user does not exists"
    exit 1
  fi
}


function main() {
  check_user
  install_pkg openssh-server

  echo "vermin ALL=(ALL) NOPASSWD:ALL" | sudo tee -a /etc/sudoers
  sudo mkdir -p /home/vermin/.ssh &&
    wget https://raw.githubusercontent.com/mhewedy/vermin/master/etc/keys/vermin_rsa.pub -O - | sudo tee -a /home/vermin/.ssh/authorized_keys &&
    sudo chmod 400 /home/vermin/.ssh/authorized_keys && sudo chown vermin:vermin /home/vermin/.ssh/authorized_keys

  # install virtualbox guest utils
  #sudo apt-get install virtualbox-guest-utils
  #sudo adduser $USER vboxsf
}

main
