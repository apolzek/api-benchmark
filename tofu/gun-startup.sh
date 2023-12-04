#!/bin/bash
echo "SERVER_API_IP=${SERVER_API_IP}" >> /etc/environment

sudo apt-get update
sudo apt-get install -y curl 

# Installing vegeta
VEGETA_VERSION=
curl -Lo vegeta.tar.gz "https://github.com/tsenart/vegeta/releases/latest/download/vegeta_$(curl -s "https://api.github.com/repos/tsenart/vegeta/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')_linux_amd64.tar.gz"

mkdir vegeta-temp
tar xf vegeta.tar.gz -C vegeta-temp

sudo mv vegeta-temp/vegeta /usr/local/bin

rm -rf vegeta.tar.gz
rm -rf vegeta-temp
