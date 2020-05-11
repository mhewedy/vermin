sudo hostnamectl set-hostname "$(hostname -I | awk '{print $1}')"

echo "export PS1=\"\[\e[1;35m\]\u\[\033[m\]@\[\e[1;92m\]$(hostname -I | awk '{print $1}')\[\033[m\]:\w \$ \"" >>~/.bashrc

## Fix IP Addr

######################################################
# https://github.com/geerlingguy/packer-boxes/issues/7
# shellcheck disable=SC2032
function wait_for_apt_lock() {
  while [ "" = "" ]; do
    eval "$1" 2>/dev/null
    if [ $? -eq 0 ]; then
      break
    fi
    sleep 3
    echo "Waiting for apt lock ..."
  done
}
######################################################

sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys CC86BB64
sudo add-apt-repository ppa:rmescandon/yq -y
wait_for_apt_lock "sudo apt-get update -y"
wait_for_apt_lock "sudo apt-get install yq -y"
sudo yq w -i /etc/netplan/00-installer-config.yaml network.ethernets.enp0s3.dhcp4 false
sudo yq w -i /etc/netplan/00-installer-config.yaml network.ethernets.enp0s3.addresses[+] "$(hostname -I | awk '{print $1}')/24"
sudo yq w -i /etc/netplan/00-installer-config.yaml network.ethernets.enp0s3.gateway4 192.168.100.1
sudo yq w -i /etc/netplan/00-installer-config.yaml network.ethernets.enp0s3.nameservers.addresses[+] 8.8.8.8
sudo yq w -i /etc/netplan/00-installer-config.yaml network.ethernets.enp0s3.nameservers.addresses[+] 8.8.4.4

sudo netplan apply
