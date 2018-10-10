#!/bin/bash

sudo yum install -y git mod_ssl policycoreutils-python
sudo yum update -y nss curl libcurl

#download go
which go 2> /dev/null
if [ $? -ne 0 ]; then 
    wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.11.1.linux-amd64.tar.gz
    echo 'Add go binary into your path PATH:$PATH:/usr/local/go/bin'
fi

port=1443
localPort=1444

echo "Adding exception for port ${port} to iptables"
iptables --list -n | grep $port > /dev/null
if [ $? -ne 0 ]; then
    iptables -I INPUT 2 -p tcp -m state --state NEW -m tcp --dport $port -j ACCEPT
    sudo service iptables save
    sudo service iptables restart
else
    echo "Port ${port} is already present"
fi

#enable port 1443 through SeLinux
echo "Adding exception for port ${port} to SeLinux"
sudo semanage port -a -t http_port_t -p tcp $port
sudo semanage port -a -t http_port_t -p tcp $localPort
