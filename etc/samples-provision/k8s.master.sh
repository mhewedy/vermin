# Set hostname and disable swap

sudo hostnamectl set-hostname "$(hostname -I | awk '{print $1}')"
sudo swapoff -a && sudo sed -i '/ swap / s/^/#/' /etc/fstab

## Fix IP Addr

sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys CC86BB64
sudo add-apt-repository ppa:rmescandon/yq
sudo apt update
sudo apt install yq -y

sudo yq w -i /etc/netplan/50-cloud-init.yaml network.ethernets.enp0s3.dhcp4 false
sudo yq w -i /etc/netplan/50-cloud-init.yaml network.ethernets.enp0s3.addresses[+] "$(hostname -I | awk '{print $1}')/24"
sudo yq w -i /etc/netplan/50-cloud-init.yaml network.ethernets.enp0s3.gateway4 192.168.100.1
sudo yq w -i /etc/netplan/50-cloud-init.yaml network.ethernets.enp0s3.nameservers.addresses[+] 8.8.8.8
sudo yq w -i /etc/netplan/50-cloud-init.yaml network.ethernets.enp0s3.nameservers.addresses[+] 8.8.4.4

sudo netplan apply

############# start k8s installation

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
cat << EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF

###

sudo apt-get update -y
sudo apt-get install -y docker-ce=18.06.1~ce~3-0~ubuntu kubelet=1.13.5-00 kubeadm=1.13.5-00 kubectl=1.13.5-00
sudo apt-mark hold docker-ce kubelet kubeadm kubectl

###

echo "net.bridge.bridge-nf-call-iptables=1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

###

sudo kubeadm init --pod-network-cidr=10.244.0.0/16 &> ~/kubeadm.log
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown "$(id -u):$(id -g)" $HOME/.kube/config
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/bc79dd1505b0c8681ece4de4c0d86c5cd2643275/Documentation/kube-flannel.yml


