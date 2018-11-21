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

# load mysqld on startup
sudo service mysqld start
sudo chkconfig --level 345 mysqld on

localPort=1444

#enable port 1443 through SeLinux
echo "Adding exception for port ${localPort} to SeLinux"
sudo semanage port -a -t http_port_t -p tcp $localPort
