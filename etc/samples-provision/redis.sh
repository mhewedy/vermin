# tested on Ubuntu focal
sudo apt update -y
sudo apt install redis-server -y
sudo sed -i 's/^supervised no/supervised systemd/' /etc/redis/redis.conf
sudo systemctl restart redis.service
sudo systemctl status redis

