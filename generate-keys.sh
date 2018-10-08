#!/bin/bash

if [ "$CERTS_PATH" == "" ]
then
    echo "CERTS_PATH is not set!!!"
    exit 1
fi

local=./keys

echo "Key and certificate will be saved under ${local} directory with server.key, server.crt respectively."

mkdir $local 2> /dev/null

openssl genrsa -out ./keys/server.key 2048
openssl req -new -x509 -sha256 -subj "/C=CA/ST=ON/L=Toronto/O=High End Security Inc./OU=BTN710" -key ./keys/server.key -out ./keys/server.crt -days 365

certs=${CERTS_PATH}/.certs

echo "Copy ${local}/* inside ${certs} for Apache to access."

sudo mkdir $certs 2> /dev/null

sudo cp $local/* $certs

echo "Transfering Apache config file"

#Replace "REPLACE_ME" string with the path to certs
replace_string="REPLACE_ME"

awk "{gsub(\"$replace_string\", \"$certs\")}1" ./go-server.conf > ./keys/go-server.conf

sudo cp ./keys/go-server.conf /etc/httpd/conf.d

echo "Setting up database..."
read -p "Enter mysql login user: " loginUser
read -p "Enter login password: " loginPassword
read -p "Enter username to be created: " createUser
read -p "Enter password for the new user: " createPassword
read -p "Enter database to be created: " database
read -p "Enter host: " host

mysql -u $loginUser --password=$loginPassword -h $host -e "create database ${database};"

cat ./res/*.sql | mysql -u $loginUser --password=$loginPassword -h $host $database

mysql -u $loginUser --password=$loginPassword -h $host -e "
create table ${database}.users ( username varchar(50) primary key, password blob(64) not null, salt blob(32) not null );
create user ${createUser}@'%' identified by '${createPassword}';
grant select on ${database}.users to ${createUser};
grant insert on ${database}.users to ${createUser};
grant select on ${database}.system_access_data to ${createUser};
grant select on ${database}.staff_data to ${createUser};
"

