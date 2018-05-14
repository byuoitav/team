#!/bin/sh

VERSION=17.09.0

sudo apt-get autoremove -y docker-ce \
    && sudo apt-get purge docker-engine -y \
    && sudo rm -rf /etc/docker/ \
    && sudo rm -f /etc/systemd/system/multi-user.target.wants/docker.service \
    && sudo rm -rf /var/lib/docker \
    &&  sudo systemctl daemon-reload

sudo apt update \
    sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg2 \
    software-properties-common

curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -

echo "deb [arch=armhf] https://download.docker.com/linux/debian \
    $(lsb_release -cs) stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list
sudo apt update

sudo apt-cache policy docker-ce

sudo apt install -y docker-ce=$VERSION~ce-0~debian

sudo apt-mark hold docker-ce
sudo printf "docker-ce hold\n" | sudo dpkg --set-selections

~/update_docker_containers.sh
