#!/usr/bin/env bash

go build

cp -f ./client cli1/
cp -f ./client cli2/
cp -f ./client cli3/
cp -f ./client cli4/
cp -f ./client cli5/

sed "s/\"address\": \"127.0.0.1:1234\"/\"address\": \"127.0.0.1:1234\"/g" ./client_conf.json > ./cli1/client_conf.json
sed "s/\"address\": \"127.0.0.1:1234\"/\"address\": \"127.0.0.1:1235\"/g" ./client_conf.json > ./cli2/client_conf.json
sed "s/\"address\": \"127.0.0.1:1234\"/\"address\": \"127.0.0.1:1236\"/g" ./client_conf.json > ./cli3/client_conf.json
sed "s/\"address\": \"127.0.0.1:1234\"/\"address\": \"127.0.0.1:1237\"/g" ./client_conf.json > ./cli4/client_conf.json
sed "s/\"address\": \"127.0.0.1:1234\"/\"address\": \"127.0.0.1:1238\"/g" ./client_conf.json > ./cli5/client_conf.json

