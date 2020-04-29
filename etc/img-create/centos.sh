sudo yum install -y openssh-server
# sudo useradd -m -s /bin/bash vermin   # vm should have a vermin/vermin user/pass
# echo "vermin:vermin" | sudo chpasswd
echo "vermin ALL=(ALL) NOPASSWD:ALL" | sudo tee -a /etc/sudoers
sudo mkdir -p /home/vermin/.ssh &&
  wget https://raw.githubusercontent.com/mhewedy/vermin/master/etc/keys/vermin_rsa.pub -O - | sudo tee -a /home/vermin/.ssh/authorized_keys &&
  sudo chmod 400 /home/vermin/.ssh/authorized_keys && sudo chown vermin:vermin /home/vermin/.ssh/authorized_keys
