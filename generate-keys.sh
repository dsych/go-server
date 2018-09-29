#!/bin/bash

echo "Key and certificate will be saved under ./keys directory with server.key, server.crt respectively"

mkdir ./keys 2> /dev/null

openssl genrsa -out ./keys/server.key 2048
openssl req -new -x509 -sha256 -key ./keys/server.key -out ./keys/server.crt -days 365
