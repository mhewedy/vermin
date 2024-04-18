# HW requirements: 1 CPUS 2G RAM
# OS: Ubuntu focal
set -ex

# Set hostname and disable swap
sudo hostnamectl set-hostname "$(hostname -I | awk '{print $1}')"
echo "export PS1=\"\[\e[1;35m\]\u\[\033[m\]@\[\e[1;92m\]$(hostname -I | awk '{print $1}')\[\033[m\]:\w \$ \"" >>~/.bashrc

sudo swapoff -a && sudo sed -i 's/\/swap/#\/swap/g' /etc/fstab

## Fix IP Addr
sudo add-apt-repository ppa:rmescandon/yq -y
sudo apt-get update -y
sudo apt-get install yq -y

# Create a new empty file (assuming the desired filename is 00-installer-config.yaml)
sudo touch /etc/netplan/00-installer-config.yaml

sudo yq eval '.network.ethernets.enp0s3.dhcp4 = false' -i /etc/netplan/00-installer-config.yaml
sudo yq eval ".network.ethernets.enp0s3.addresses += [\"$(hostname -I | awk '{print $1}')/24\"]" -i /etc/netplan/00-installer-config.yaml
sudo yq eval ".network.ethernets.enp0s3.gateway4 = \"$(hostname -I | cut -d '.' -f 1-3).1\"" -i /etc/netplan/00-installer-config.yaml
sudo yq eval '.network.ethernets.enp0s3.nameservers.addresses += ["8.8.8.8"]' -i /etc/netplan/00-installer-config.yaml
sudo yq eval '.network.ethernets.enp0s3.nameservers.addresses += ["8.8.4.4"]' -i /etc/netplan/00-installer-config.yaml

sudo netplan apply
sleep 1

############# install docker & k8s

sudo apt-get update && sudo apt-get install -y apt-transport-https curl

sudo mkdir -p -m 755 /etc/apt/keyrings


# Add Kubernetes APT repository

curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

## install docker, kubelet, kubeadm and kubectl

# Update to the latest Docker version
sudo apt-get update
sudo apt-get install -y docker.io

# Update to the latest Kubernetes version
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

echo "net.bridge.bridge-nf-call-iptables=1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

## configure docker

cat <<EOF | sudo tee /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

sudo mkdir -p /etc/systemd/system/docker.service.d

sudo systemctl daemon-reload
sudo systemctl restart docker
