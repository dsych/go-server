#!/bin/bash

localDirectory="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

local="${localDirectory}/keys"

mkdir $local 2> /dev/null


echo "Session private key will be stored under ${local}/session.key."

openssl genrsa -out $local/session.key 2048

echo "Key and certificate will be saved under ${local} directory with server.key, server.crt respectively."

mkdir $local 2> /dev/null

openssl genrsa -out $local/server.key 2048
openssl req -new -x509 -sha256 -subj "/C=CA/ST=ON/L=Toronto/O=High End Security Inc./OU=BTN710" -key $local/server.key -out $local/server.crt -days 365

echo "Transfering Apache config file"

sudo cp $localDirectory/go-server.conf /etc/httpd/conf.d
sudo cp $local/server.crt /etc/pki/tls/private
sudo cp $local/server.key /etc/pki/tls/private
sudo cp $local/server.crt /etc/pki/tls/certs

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

