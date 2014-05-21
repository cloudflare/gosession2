#!/bin/bash -xe
rm -rf pki
./easyrsa init-pki
./easyrsa build-ca
./easyrsa build-client-full client
./easyrsa build-server-full server
cp -v pki/ca.crt .
cp -v pki/issued/*.crt .
openssl rsa -in ./pki/private/client.key -out client.key
openssl rsa -in ./pki/private/server.key -out server.key
